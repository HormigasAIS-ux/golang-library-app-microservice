package repository

import (
	author_grpc "auth_service/interface/grpc/genproto/author"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorRepo struct {
	authorGrpcServiceConn author_grpc.AuthorServiceClient
}

type IAuthorRepo interface {
	RpcCreateAuthor(
		ctx context.Context,
		payload *author_grpc.CreateAuthorReq,
	) (*author_grpc.CreateAuthorResp, codes.Code, error)
}

func NewAuthorRepo(authorGrpcServiceConn author_grpc.AuthorServiceClient) IAuthorRepo {
	return &AuthorRepo{
		authorGrpcServiceConn: authorGrpcServiceConn,
	}
}

func (repo *AuthorRepo) RpcCreateAuthor(
	ctx context.Context,
	payload *author_grpc.CreateAuthorReq,
) (*author_grpc.CreateAuthorResp, codes.Code, error) {
	res, err := repo.authorGrpcServiceConn.CreateAuthor(ctx, payload)
	code := status.Code(err)
	return res, code, err
}
