package http_response

import (
	"auth_service/domain/dto"
	error_utils "auth_service/utils/error"

	"github.com/gin-gonic/gin"
)

type HttpResponseWriter struct{}

type IHttpResponseWriter interface {
	HTTPCustomErr(ctx *gin.Context, err error)
	HTTPJson(ctx *gin.Context, code int, message string, detail string, data interface{})
	HTTPJsonOK(ctx *gin.Context, data interface{})
}

func NewHttpResponseWriter() IHttpResponseWriter {
	return &HttpResponseWriter{}
}

func (r *HttpResponseWriter) HTTPCustomErr(ctx *gin.Context, err error) {
	customErr, ok := err.(*error_utils.CustomErr)
	if ok {
		ctx.JSON(customErr.HttpCode, dto.BaseJSONResp{
			Code:    customErr.HttpCode,
			Message: customErr.Message,
			Detail:  customErr.Detail,
			Data:    customErr.Data,
		})
		return
	}
	ctx.JSON(500, dto.BaseJSONResp{
		Code:    500,
		Message: "internal server error",
		Detail:  err.Error(),
		Data:    nil,
	})
}

func (r *HttpResponseWriter) HTTPJson(ctx *gin.Context, code int, message string, detail string, data interface{}) {
	ctx.JSON(code, dto.BaseJSONResp{
		Code:    code,
		Message: message,
		Detail:  detail,
		Data:    data,
	})
}

func (r *HttpResponseWriter) HTTPJsonOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, dto.BaseJSONResp{
		Code:    200,
		Message: "OK",
		Detail:  "",
		Data:    data,
	})
}
