package handler

import (
	"auth_service/domain/dto"
	ucase "auth_service/usecase"
	"auth_service/utils/http_response"

	"github.com/gin-gonic/gin"
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

// Register
// @Summary register new user
// @Tags Auth
// @Success 200 {object} dto.BaseJSONResp{data=dto.RegisterUserRespData}
// @Router /auth/register [post]
// @param payload  body  dto.RegisterUserReq  true "payload"
func (h *AuthHandler) Register(ctx *gin.Context) {
	var payload dto.RegisterUserReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	data, err := h.authUcase.Register(payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// Login
// @Summary login
// @Tags Auth
// @Success 200 {object} dto.BaseJSONResp{data=dto.LoginRespData}
// @Router /auth/login [post]
// @param payload  body  dto.LoginReq  true "payload"
func (h *AuthHandler) Login(ctx *gin.Context) {
	var payload dto.LoginReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	data, err := h.authUcase.Login(payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// Refresh Token
// @Tags Auth
// @Router /auth/refresh-token [post]
// @param payload  body  dto.RefreshTokenReq  true "payload"
// @Success 200 {object} dto.BaseJSONResp{data=dto.RefreshTokenRespData}
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var payload dto.RefreshTokenReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	data, err := h.authUcase.RefreshToken(payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// Check Token
// @Tags Auth
// @Router /auth/check-token [post]
// @param payload  body  dto.CheckTokenReq  true "payload"
// @Success 200 {object} dto.BaseJSONResp{data=dto.CheckTokenRespData}
func (h *AuthHandler) CheckToken(ctx *gin.Context) {
	var payload dto.CheckTokenReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	data, err := h.authUcase.CheckToken(payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}
