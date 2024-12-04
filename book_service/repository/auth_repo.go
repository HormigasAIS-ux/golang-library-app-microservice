package repository

import (
	pb "book_service/interface/grpc/genproto/auth"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthRepo struct {
	authGrpcServiceConn pb.AuthServiceClient
}

type IAuthRepo interface {
	RpcGetUserByUUID(
		ctx context.Context,
		payload *pb.GetUserByUUIDRequest,
	) (*pb.GetUserByUUIDResponse, codes.Code, error)
	RpcCreateUser(
		ctx context.Context,
		payload *pb.CreateUserReq,
	) (*pb.CreateUserResp, codes.Code, error)
	RpcUpdateUser(
		ctx context.Context,
		payload *pb.UpdateUserReq,
	) (*pb.UpdateUserResp, codes.Code, error)
	RpcDeleteUser(
		ctx context.Context,
		payload *pb.DeleteUserReq,
	) (*pb.DeleteUserResp, codes.Code, error)
}

func NewAuthRepo(authGrpcServiceConn pb.AuthServiceClient) IAuthRepo {
	return &AuthRepo{
		authGrpcServiceConn: authGrpcServiceConn,
	}
}

func (repo *AuthRepo) RpcGetUserByUUID(
	ctx context.Context,
	payload *pb.GetUserByUUIDRequest,
) (*pb.GetUserByUUIDResponse, codes.Code, error) {
	res, err := repo.authGrpcServiceConn.GetUserByUUID(ctx, payload)
	code := status.Code(err)
	return res, code, err
}

func (repo *AuthRepo) RpcCreateUser(
	ctx context.Context,
	payload *pb.CreateUserReq,
) (*pb.CreateUserResp, codes.Code, error) {
	res, err := repo.authGrpcServiceConn.CreateUser(ctx, payload)
	code := status.Code(err)
	return res, code, err
}

func (repo *AuthRepo) RpcUpdateUser(
	ctx context.Context,
	payload *pb.UpdateUserReq,
) (*pb.UpdateUserResp, codes.Code, error) {
	res, err := repo.authGrpcServiceConn.UpdateUser(ctx, payload)
	code := status.Code(err)
	return res, code, err
}

func (repo *AuthRepo) RpcDeleteUser(
	ctx context.Context,
	payload *pb.DeleteUserReq,
) (*pb.DeleteUserResp, codes.Code, error) {
	res, err := repo.authGrpcServiceConn.DeleteUser(ctx, payload)
	code := status.Code(err)
	return res, code, err
}
