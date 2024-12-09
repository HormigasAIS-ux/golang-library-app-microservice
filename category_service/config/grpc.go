package config

import (
	book_grpc "category_service/interface/grpc/genproto/book"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// func NewAuthGrpcServiceClient() auth_pb.AuthServiceClient {
// 	conn, err := grpc.NewClient(Envs.AUTH_GRPC_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		logger.Fatalf("Failed to connect to auth grpc service: %v", err)
// 	}
// 	authServiceClient := auth_pb.NewAuthServiceClient(conn)
// 	return authServiceClient
// }

// func NewAuthorGrpcServiceClient() author_grpc.AuthorServiceClient {
// 	conn, err := grpc.NewClient(Envs.AUTHOR_GRPC_SERVCICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		logger.Fatalf("Failed to connect to auth grpc service: %v", err)
// 	}
// 	authServiceClient := author_grpc.NewAuthorServiceClient(conn)
// 	return authServiceClient
// }

func NewBookGrpcServiceClient() book_grpc.BookServiceClient {
	conn, err := grpc.NewClient(Envs.AUTHOR_GRPC_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to connect to auth grpc service: %v", err)
	}
	authServiceClient := book_grpc.NewBookServiceClient(conn)
	return authServiceClient
}
