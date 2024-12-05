package error_utils

import "google.golang.org/grpc/codes"

type CustomErr struct {
	HttpCode int
	GrpcCode codes.Code
	Message  string
	Detail   interface{}
	Data     interface{}
}

func (slf *CustomErr) Error() string {
	if slf.Detail != "" || slf.Detail != nil {
		err, ok := slf.Detail.(error)
		if ok {
			return err.Error()
		}

		detail, ok := slf.Detail.(string)
		if ok {
			return detail
		}

		return "unknown error"
	}
	return slf.Message
}
