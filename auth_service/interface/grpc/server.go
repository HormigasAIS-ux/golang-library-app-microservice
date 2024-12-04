package grpc

import (
	"auth_service/config"
	interface_pkg "auth_service/interface"
	"auth_service/interface/grpc/genproto/auth"
	"auth_service/interface/grpc/handler"
	"fmt"
	"log"
	"net"

	"github.com/op/go-logging"
	"google.golang.org/grpc"
)

var logger = logging.MustGetLogger("main")

func SetupServer(commonDependencies interface_pkg.CommonDependency) {
	// setup listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Envs.GRPC_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// new grpc server
	grpcServer := grpc.NewServer()

	// register service handler
	logger.Debugf("ucase: %v", commonDependencies.AuthUcase)
	authServiceHandler := handler.NewAuthServiceHandler(commonDependencies.AuthUcase, commonDependencies.UserUcase)
	auth.RegisterAuthServiceServer(grpcServer, authServiceHandler)

	// Start the server
	fmt.Printf("Starting gRPC server on port :%v...", config.Envs.GRPC_PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
