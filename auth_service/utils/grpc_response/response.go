package grpc_response

import (
	"auth_service/domain/dto"
	error_utils "auth_service/utils/error"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcResponseWriter struct{}

type IGrpcResponseWriter interface {
	RespCustomError(err error) (interface{}, error)
}

func NewGrpcResponseWriter() IGrpcResponseWriter {
	return &GrpcResponseWriter{}
}

func (r *GrpcResponseWriter) RespCustomError(err error) (interface{}, error) {
	customErr, ok := err.(*error_utils.CustomErr)

	if ok {
		// handle emtpy grpc code
		code := customErr.GrpcCode
		if code == 0 {
			code = codes.Internal
		}

		return dto.BaseGrpcResp{
			Code:    code,
			Message: customErr.Error(),
			Detail:  "",
			Data:    nil,
		}, status.Errorf(customErr.GrpcCode, customErr.Error())
	}
	return dto.BaseGrpcResp{
		Code:    codes.Internal,
		Message: "internal server error",
		Detail:  err.Error(),
		Data:    nil,
	}, status.Errorf(codes.Internal, err.Error())
}
