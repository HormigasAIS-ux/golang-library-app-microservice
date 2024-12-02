package main

import (
	"auth_service/api/rest"
	"auth_service/config"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

func init() {
	config.InitEnv("./.env")
	config.ConfigureLogger()
}

var logger = logging.MustGetLogger("main")

func runRestServer() {
	logger.Info("starting rest server...")
	router := gin.Default()
	rest.SetupServer(router)
	router.Run(fmt.Sprintf("%s:%d", config.Envs.HOST, config.Envs.PORT))
}

// @title Auth Service RESTful API
// @securitydefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme (add 'Bearer ' prefix).
func main() {
	args := os.Args
	if len(args) == 1 { // run as a rest server
		runRestServer()
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
					runRestServer()
				case "grpc":
					logger.Fatalf("grpc not implemented yet")
				default:
					logger.Fatalf("invalid argument: %s", arg)
				}
			}
		}
	}
}
