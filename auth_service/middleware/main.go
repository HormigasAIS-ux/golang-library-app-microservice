package middleware

import (
	ucase "auth_service/usecase"
	"auth_service/utils/http_response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(respWriter http_response.IResponseWriter, authService ucase.IAuthUcase) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = c.GetHeader("Bearer")

		c.Next()
	}
}
