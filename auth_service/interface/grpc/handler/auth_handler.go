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
	authUcase ucase.IAuthUcase
	userUcase ucase.IUserUcase
}

func NewAuthServiceHandler(authUcase ucase.IAuthUcase, userUcase ucase.IUserUcase) *AuthServiceHandler {
	return &AuthServiceHandler{authUcase: authUcase, userUcase: userUcase}
}

func (h *AuthServiceHandler) CheckToken(ctx context.Context, req *auth.CheckTokenRequest) (*auth.CheckTokenResponse, error) {
	// payload validation
	payload := dto.CheckTokenReq{
		AccessToken: req.AccessToken,
	}
	if payload.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "missing access token")
	}

	raw, err := h.authUcase.CheckToken(payload)
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
		Role:     raw.Role,
	}

	return resp, nil
}

func (h *AuthServiceHandler) GetUserByUUID(
	ctx context.Context,
	req *auth.GetUserByUUIDRequest,
) (*auth.GetUserByUUIDResponse, error) {
	// payload validation
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "missing uuid")
	}

	raw, err := h.userUcase.GetByUUID(ctx, nil, req.Uuid)
	if err != nil {
		customErr, ok := err.(*error_utils.CustomErr)
		if ok {
			return nil, status.Errorf(customErr.GrpcCode, customErr.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &auth.GetUserByUUIDResponse{
		Uuid:     raw.UUID,
		Username: raw.Username,
		Email:    raw.Email,
		Role:     raw.Role,
	}

	return resp, nil
}

func (h *AuthServiceHandler) CreateUser(
	ctx context.Context,
	req *auth.CreateUserReq,
) (*auth.CreateUserResp, error) {
	// payload validation
	if req.Username == "" ||
		req.Email == "" ||
		req.Password == "" ||
		req.Role == "" {
		return nil, status.Error(codes.InvalidArgument, "fields required")
	}

	raw, err := h.userUcase.CreateUser(ctx, nil, dto.CreateUserReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		customErr, ok := err.(*error_utils.CustomErr)
		if ok {
			return nil, status.Errorf(customErr.GrpcCode, customErr.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &auth.CreateUserResp{
		Uuid:     raw.UUID,
		Username: raw.Username,
		Email:    raw.Email,
		Role:     raw.Role,
	}

	return resp, nil
}

func (h *AuthServiceHandler) UpdateUser(
	ctx context.Context,
	req *auth.UpdateUserReq,
) (*auth.UpdateUserResp, error) {
	// payload validation
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "missing uuid")
	}

	// convert any empty string to nil
	dtoPayload := dto.UpdateUserReq{}
	if req.Username == "" {
		dtoPayload.Username = nil
	} else {
		dtoPayload.Username = &req.Username
	}
	if req.Email == "" {
		dtoPayload.Email = nil
	} else {
		dtoPayload.Email = &req.Email
	}
	if req.Password == "" {
		dtoPayload.Password = nil
	} else {
		dtoPayload.Password = &req.Password
	}
	if req.Role == "" {
		dtoPayload.Role = nil
	} else {
		dtoPayload.Role = &req.Role
	}

	raw, err := h.userUcase.UpdateUser(ctx, nil, req.Uuid, dtoPayload)
	if err != nil {
		customErr, ok := err.(*error_utils.CustomErr)
		if ok {
			return nil, status.Errorf(customErr.GrpcCode, customErr.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &auth.UpdateUserResp{
		Uuid:     raw.UUID,
		Username: raw.Username,
		Email:    raw.Email,
		Role:     raw.Role,
	}

	return resp, nil
}

func (h *AuthServiceHandler) DeleteUser(
	ctx context.Context,
	req *auth.DeleteUserReq,
) (*auth.DeleteUserResp, error) {
	// payload validation
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "missing uuid")
	}

	raw, err := h.userUcase.DeleteUser(ctx, nil, req.Uuid)
	if err != nil {
		customErr, ok := err.(*error_utils.CustomErr)
		if ok {
			return nil, status.Errorf(customErr.GrpcCode, customErr.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &auth.DeleteUserResp{
		Uuid:     raw.UUID,
		Username: raw.Username,
		Email:    raw.Email,
		Role:     raw.Role,
	}

	return resp, nil
}
