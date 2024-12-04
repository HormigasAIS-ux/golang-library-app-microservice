package dto

import (
	"google.golang.org/grpc/codes"
)

type BaseJSONResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
	Data    interface{} `json:"data"`
}

type BaseGrpcResp struct {
	Code    codes.Code  `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
	Data    interface{} `json:"data"`
}
