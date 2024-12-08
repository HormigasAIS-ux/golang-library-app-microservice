package config

import (
	auth_grpc "author_service/interface/grpc/genproto/auth"
	book_grpc "author_service/interface/grpc/genproto/book"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthGrpcServiceClient() auth_grpc.AuthServiceClient {
	conn, err := grpc.NewClient(Envs.AUTH_GRPC_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to connect to auth grpc service: %v", err)
	}
	authServiceClient := auth_grpc.NewAuthServiceClient(conn)
	return authServiceClient
}

func NewBookGrpcServiceClient() book_grpc.BookServiceClient {
	conn, err := grpc.NewClient(Envs.BOOK_GRPC_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to connect to book grpc service: %v", err)
	}
	bookServiceClient := book_grpc.NewBookServiceClient(conn)
	return bookServiceClient
}
