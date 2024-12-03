package main

import (
	"auth_service/config"
	"auth_service/domain/model"
	interface_pkg "auth_service/interface"
	"auth_service/interface/grpc"
	"auth_service/interface/rest"
	"auth_service/repository"
	ucase "auth_service/usecase"
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
	gormDB := config.NewPostgresqlDB()

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

	// ucases
	authUcase := ucase.NewAuthUcase(userRepo, refreshTokenRepo)

	dependencies := interface_pkg.CommonDependency{
		AuthUcase: authUcase,
	}

	// seed data
	err = seeder_util.SeedUser(userRepo)
	if err != nil {
		logger.Fatalf("failed to seed user: %v", err)
	}

	args := os.Args
	if len(args) == 1 { // run as a rest server
		logger.Info("starting rest server...")
		rest.SetupServer(dependencies)
	} else if len(args) > 1 {
		validArgVariables := []string{"server"}

		// validate args
		for _, arg := range args[1:] {
			logger.Debugf("arg: %s", arg)
			for _, validArgVariable := range validArgVariables {
				if strings.Contains(arg, fmt.Sprintf("--%s=", validArgVariable)) {
					continue
				}
				logger.Fatalf("invalid argument: %s", arg)
			}
		}

		// process args
		for _, arg := range args[1:] {
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
