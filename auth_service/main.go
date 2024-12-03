package main

import (
	"auth_service/config"
	"auth_service/interface/grpc"
	"auth_service/interface/rest"
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
	args := os.Args
	if len(args) == 1 { // run as a rest server
		logger.Info("starting rest server...")
		rest.SetupServer()
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
					rest.SetupServer()
				case "grpc":
					logger.Info("starting grpc server...")
					grpc.SetupServer()
				default:
					logger.Fatalf("invalid argument: %s", arg)
				}
			}
		}
	}
}
