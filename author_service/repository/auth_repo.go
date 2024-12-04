package repository

import (
	pb "author_service/interface/grpc/genproto/auth"
)

type AuthRepo struct {
	authGrpcServiceConn pb.AuthServiceClient
}

type IAuthRepo interface {
}

func NewAuthRepo(authGrpcServiceConn pb.AuthServiceClient) IAuthRepo {
	return &AuthRepo{
		authGrpcServiceConn: authGrpcServiceConn,
	}
}
