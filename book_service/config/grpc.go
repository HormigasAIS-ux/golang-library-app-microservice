package config

import (
	author_pb "book_service/interface/grpc/genproto/author"

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

func NewAuthorGrpcServiceClient() author_pb.AuthorServiceClient {
	conn, err := grpc.NewClient(Envs.AUTHOR_GRPC_SERVCICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to connect to auth grpc service: %v", err)
	}
	authServiceClient := author_pb.NewAuthorServiceClient(conn)
	return authServiceClient
}
