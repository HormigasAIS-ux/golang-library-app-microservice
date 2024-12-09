package main

import (
	"category_service/config"
	"category_service/domain/model"
	interface_pkg "category_service/interface"
	"category_service/interface/grpc"
	"category_service/interface/rest"
	"category_service/repository"
	ucase "category_service/usecase"
	"category_service/utils/helper"
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

// @title Category Service RESTful API
// @securitydefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme (add 'Bearer ' prefix).
func main() {
	logger.Debugf("Envs: %v", helper.PrettyJson(config.Envs))
	gormDB := config.NewPostgresqlDB()
	bookGrpcServiceClient := config.NewBookGrpcServiceClient()
	// authGrpcServiceClient := config.NewAuthGrpcServiceClient()
	// authorGrpcServiceClient := config.NewAuthorGrpcServiceClient()

	// migrations
	err := gormDB.AutoMigrate(
		&model.Category{},
	)
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	// repositories
	categoryRepo := repository.NewCategoryRepo(gormDB)

	// ucases
	categoryUcase := ucase.NewCategoryUcase(categoryRepo, bookGrpcServiceClient)
	dependencies := interface_pkg.CommonDependency{
		CategoryUcase: categoryUcase,
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
		variables = append(variables, validPreRunArgVariables...)

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
		variables = append(variables, postArgs...)
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
