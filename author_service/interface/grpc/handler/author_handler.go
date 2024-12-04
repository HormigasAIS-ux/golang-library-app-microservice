package handler

import (
	author_pb "author_service/interface/grpc/genproto/author"
	ucase "author_service/usecase"
)

type AuthorServiceHandler struct {
	author_pb.UnimplementedAuthServiceServer
	authorService ucase.IAuthorUcase
}

func NewAuthorServiceHandler(authService ucase.IAuthorUcase) *AuthorServiceHandler {
	return &AuthorServiceHandler{authorService: authService}
}
