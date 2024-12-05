package repository

import (
	author_pb "book_service/interface/grpc/genproto/author"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorRepo struct {
	authorGrpcServiceConn author_pb.AuthorServiceClient
}

type IAuthorRepo interface {
	RpcCreateAuthor(
		ctx context.Context,
		payload *author_pb.CreateAuthorReq,
	) (*author_pb.CreateAuthorResp, codes.Code, error)
}

func NewAuthorRepo(authorGrpcServiceConn author_pb.AuthorServiceClient) IAuthorRepo {
	return &AuthorRepo{
		authorGrpcServiceConn: authorGrpcServiceConn,
	}
}
func (repo *AuthorRepo) RpcCreateAuthor(
	ctx context.Context,
	payload *author_pb.CreateAuthorReq,
) (*author_pb.CreateAuthorResp, codes.Code, error) {
	res, err := repo.authorGrpcServiceConn.CreateAuthor(ctx, payload)
	code := status.Code(err)
	return res, code, err
}

func (repo *AuthorRepo) RpcGetAuthorByUserUUID(
	ctx context.Context,
	payload *author_pb.GetAuthorByUserUUIDReq,
) (*author_pb.GetAuthorByUserUUIDResp, codes.Code, error) {
	res, err := repo.authorGrpcServiceConn.GetAuthorByUserUUID(ctx, payload)
	code := status.Code(err)
	return res, code, err
}
