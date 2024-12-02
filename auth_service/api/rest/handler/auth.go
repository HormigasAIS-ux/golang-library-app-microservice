package handler

import (
	ucase "auth_service/usecase"
	"auth_service/utils/http_response"
)

type AuthHandler struct {
	respWriter http_response.IResponseWriter
	authUcase  ucase.IAuthUcase
}

type IAuthHandler interface {
}

func NewAuthHandler(respWriter http_response.IResponseWriter, authUcase ucase.IAuthUcase) AuthHandler {
	return AuthHandler{
		respWriter: respWriter,
		authUcase:  authUcase,
	}
}
