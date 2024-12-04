package config

import (
	pb "author_service/interface/grpc/genproto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthGrpcServiceClient() pb.AuthServiceClient {
	conn, err := grpc.NewClient(Envs.AUTH_GRPC_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to connect to auth grpc service: %v", err)
	}
	authServiceClient := pb.NewAuthServiceClient(conn)
	return authServiceClient
}
