package error_utils

import "google.golang.org/grpc/codes"

type CustomErr struct {
	HttpCode int
	GrpcCode codes.Code
	Message  string
	Detail   string
	Data     interface{}
}

func (slf *CustomErr) Error() string {
	if slf.Detail != "" {
		return slf.Detail
	}
	return slf.Message
}
