package main

import (
	"book_service/config"
	"book_service/domain/model"
	interface_pkg "book_service/interface"
	"book_service/interface/grpc"
	"book_service/interface/rest"
	"book_service/repository"
	ucase "book_service/usecase"
	"book_service/utils/helper"
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

// @title Book Service RESTful API
// @securitydefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme (add 'Bearer ' prefix).
func main() {
	logger.Debugf("Envs: %v", helper.PrettyJson(config.Envs))
	gormDB := config.NewPostgresqlDB()
	authGrpcServiceClient := config.NewAuthGrpcServiceClient()

	// migrations
	err := gormDB.AutoMigrate(
		&model.Book{},
		&model.BookBorrow{},
	)
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	// repositories
	authRepo := repository.NewAuthRepo(authGrpcServiceClient)
	authorRepo := repository.NewAuthorRepo(gormDB)

	// ucases
	authorUcase := ucase.NewAuthorUcase(authorRepo, authRepo)
	dependencies := interface_pkg.CommonDependency{
		AuthorUcase: authorUcase,
	}

	args := os.Args
	if len(args) == 1 { // run as a rest server
		logger.Info("starting rest server...")
		rest.SetupServer(dependencies)
	} else if len(args) > 1 {
		validArgVariables := []string{"server"}
		validPreRunArgVariables := []string{}

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
			}
		}
	}
}
