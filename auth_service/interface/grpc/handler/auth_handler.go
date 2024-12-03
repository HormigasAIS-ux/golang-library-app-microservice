package handler

import (
	"auth_service/domain/dto"
	"auth_service/interface/grpc/genproto/auth"
	ucase "auth_service/usecase"
	error_utils "auth_service/utils/error"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceHandler struct {
	auth.UnimplementedAuthServiceServer
	authService ucase.IAuthUcase
}

func NewAuthServiceHandler(authService ucase.IAuthUcase) *AuthServiceHandler {
	return &AuthServiceHandler{authService: authService}
}

func (h *AuthServiceHandler) CheckToken(ctx context.Context, req *auth.CheckTokenRequest) (*auth.CheckTokenResponse, error) {
	// payload validation
	payload := dto.CheckTokenReq{
		AccessToken: req.AccessToken,
	}
	if payload.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "missing access token")
	}

	raw, err := h.authService.CheckToken(payload)
	if err != nil {
		customErr, ok := err.(*error_utils.CustomErr)
		if ok {
			return nil, status.Errorf(customErr.GrpcCode, customErr.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &auth.CheckTokenResponse{
		Uuid:     raw.UUID,
		Username: raw.Username,
		Email:    raw.Email,
		Fullname: raw.Fullname,
	}

	return resp, nil
}
