package main

import (
	"auth_service/config"
	"auth_service/domain/model"
	interface_pkg "auth_service/interface"
	"auth_service/interface/grpc"
	"auth_service/interface/rest"
	"auth_service/repository"
	ucase "auth_service/usecase"
	"auth_service/utils/helper"
	seeder_util "auth_service/utils/seeder/user"
	"fmt"
	"os"
	"strings"

	"github.com/op/go-logging"
)

func init() {
	config.InitEnv("./.env")
	config.ConfigureLogger()
}

var logger = logging.MustGetLogger("main")

// @title Auth Service RESTful API
// @securitydefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme (add 'Bearer ' prefix).
func main() {
	logger.Debugf("Envs: %v", helper.PrettyJson(config.Envs))

	gormDB := config.NewPostgresqlDB()
	authorGrpcServiceClient := config.NewAuthorGrpcServiceClient()

	// migrations
	err := gormDB.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
	)
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	// prepare dependencies
	// repositories
	userRepo := repository.NewUserRepo(gormDB)
	refreshTokenRepo := repository.NewRefreshTokenRepo(gormDB)
	authorRepo := repository.NewAuthorRepo(authorGrpcServiceClient)

	// ucases
	authUcase := ucase.NewAuthUcase(userRepo, refreshTokenRepo, authorRepo)
	userUcase := ucase.NewUserUcase(userRepo)

	dependencies := interface_pkg.CommonDependency{
		AuthUcase: authUcase,
		UserUcase: userUcase,
	}

	args := os.Args
	if len(args) == 1 { // run as a rest server
		logger.Info("starting rest server...")
		rest.SetupServer(dependencies)
	} else if len(args) > 1 {
		validArgVariables := []string{"server"}
		validPreRunArgVariables := []string{"seed"}

		// validate args
		variables := validArgVariables
		for _, preRunVariable := range validPreRunArgVariables {
			variables = append(variables, preRunVariable)
		}
		// logger.Debugf("variables: %v", variables)
		for _, arg := range args[1:] {
			valid := false
			// logger.Debugf("arg: %s", arg)
			for _, validArgVariable := range variables {
				if strings.Contains(arg, fmt.Sprintf("--%s=", validArgVariable)) {
					// logger.Debug("contains")
					valid = true
					break
				}
			}

			if !valid {
				logger.Fatalf("invalid argument: %s", arg)
			}
		}

		// group between pre variable and post variable
		preArgs := []string{}
		postArgs := []string{}
		for _, arg := range args[1:] {
			for _, preRunVariable := range validPreRunArgVariables {
				if strings.Contains(arg, fmt.Sprintf("--%s=", preRunVariable)) {
					preArgs = append(preArgs, arg)
					// logger.Debugf("preArg: %s", arg)
				}
			}

			for _, validArgVariable := range validArgVariables {
				if strings.Contains(arg, fmt.Sprintf("--%s=", validArgVariable)) {
					postArgs = append(postArgs, arg)
					// logger.Debugf("postArg: %s", arg)
				}
			}
		}

		// process args
		variables = preArgs
		for _, postArg := range postArgs {
			variables = append(variables, postArg)
		}
		for _, arg := range variables {
			if strings.Contains(arg, fmt.Sprintf("--%s=", "server")) {
				value := strings.Split(arg, "=")[1]

				switch value {
				case "rest":
					logger.Info("starting rest server...")
					rest.SetupServer(dependencies)
				case "grpc":
					logger.Info("starting grpc server...")
					grpc.SetupServer(dependencies)
				default:
					logger.Fatalf("invalid argument: %s", arg)
				}

			} else if strings.Contains(arg, fmt.Sprintf("--%s=", "seed")) {
				value := strings.Split(arg, "=")[1]

				switch value {
				case "user":
					err = seeder_util.SeedUser(userRepo, authorRepo)
					if err != nil {
						logger.Fatalf("failed to seed user: %v", err)
					}
				}
			}
		}
	}
}
