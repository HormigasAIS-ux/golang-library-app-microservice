package handler

import (
	"author_service/domain/dto"
	author_pb "author_service/interface/grpc/genproto/author"
	ucase "author_service/usecase"
	error_utils "author_service/utils/error"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorServiceHandler struct {
	author_pb.UnimplementedAuthorServiceServer
	authorService ucase.IAuthorUcase
}

func NewAuthorServiceHandler(authService ucase.IAuthorUcase) *AuthorServiceHandler {
	return &AuthorServiceHandler{authorService: authService}
}

func (r *AuthorServiceHandler) CreateAuthor(
	ctx *gin.Context,
	payload *author_pb.CreateAuthorReq,
) (*author_pb.CreateAuthorResp, error) {
	// payload validation
	payloadDto := dto.CreateNewAuthorReq{
		LastName: payload.LastName,
	}
	if payload.UserUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "user uuid is required")
	}

	if payload.FirstName == "" {
		return nil, status.Error(codes.InvalidArgument, "first name is required")
	} else {
		payloadDto.FirstName = payload.FirstName
	}

	if payload.BirthDate == "" {
		payloadDto.BirthDate = nil
	} else {
		payloadDto.BirthDate = &payload.BirthDate
	}

	if payload.Bio == "" {
		payloadDto.Bio = nil
	} else {
		payloadDto.Bio = &payload.Bio
	}

	raw, err := r.authorService.CreateNewAuthor(ctx, payloadDto)
	if err != nil {
		customErr, ok := err.(*error_utils.CustomErr)
		if ok {
			return nil, status.Errorf(customErr.GrpcCode, customErr.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &author_pb.CreateAuthorResp{
		Uuid:      raw.UUID.String(),
		UserUuid:  raw.UserUUID.String(),
		FirstName: raw.FirstName,
		LastName:  raw.LastName,
	}

	if raw.BirthDate != nil {
		resp.BirthDate = *raw.BirthDate
	} else {
		resp.BirthDate = ""
	}

	if raw.Bio != nil {
		resp.Bio = *raw.Bio
	} else {
		resp.Bio = ""
	}

	return resp, nil
}
