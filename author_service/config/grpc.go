package config

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcConnection() *grpc.ClientConn {
	conn, err := grpc.NewClient(Envs.AUTH_GRPC_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Failed to connect to auth grpc service: %v", err)
	}
	return conn
}
