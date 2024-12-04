package config

import (
	author_pb "auth_service/interface/grpc/genproto/author"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthorGrpcServiceClient() author_pb.AuthorServiceClient {
	conn, err := grpc.NewClient(Envs.AUTHOR_GRPC_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to connect to author grpc service: %v", err)
	}
	authServiceClient := author_pb.NewAuthorServiceClient(conn)
	return authServiceClient
}
