package grpc

import (
	"auth_service/config"
	"auth_service/domain/model"
	"auth_service/interface/grpc/genproto/auth"
	"auth_service/interface/grpc/handler"
	"auth_service/repository"
	ucase "auth_service/usecase"
	"fmt"
	"log"
	"net"

	"github.com/op/go-logging"
	"google.golang.org/grpc"
)

var logger = logging.MustGetLogger("main")

func SetupServer() {
	gormDB := config.NewPostgresqlDB()

	// migrations
	err := gormDB.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
	)
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	// repositories
	userRepo := repository.NewUserRepo(gormDB)
	refreshTokenRepo := repository.NewRefreshTokenRepo(gormDB)

	// ucases
	authUcase := ucase.NewAuthUcase(userRepo, refreshTokenRepo)

	// setup listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Envs.GRPC_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// new grpc server
	grpcServer := grpc.NewServer()

	// register service handler
	authServiceHandler := handler.NewAuthServiceHandler(authUcase)
	auth.RegisterAuthServiceServer(grpcServer, authServiceHandler)

	// Start the server
	fmt.Printf("Starting gRPC server on port :%v...", config.Envs.GRPC_PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
