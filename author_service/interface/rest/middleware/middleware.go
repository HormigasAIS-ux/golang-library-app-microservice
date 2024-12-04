package rest_middleware

import (
	"author_service/config"
	"author_service/domain/dto"
	"author_service/utils/http_response"
	jwt_util "author_service/utils/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(respWriter http_response.IHttpResponseWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			respWriter.HTTPJson(
				c, 401, "unauthorized", "invalid token", nil,
			)
			c.Abort()
			return
		}

		currentUser, err := jwt_util.ValidateJWT(token, config.Envs.JWT_SECRET_KEY)
		if err != nil {
			respWriter.HTTPJson(
				c, 401, "unauthorized", err.Error(), nil,
			)
			c.Abort()
			return
		}

		c.Set("currentUser", currentUser)
		c.Next()
	}
}

func AuthAdminOnlyMiddleware(respWriter http_response.IHttpResponseWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserRaw, ok := c.Get("currentUser")
		if !ok {
			respWriter.HTTPJson(
				c, 500, "internal service error", "current user not found", nil,
			)
			return
		}

		currentUser, ok := currentUserRaw.(dto.CurrentUser)
		if !ok {
			respWriter.HTTPJson(
				c, 500, "internal service error", "current user missmatched", nil,
			)
			return
		}

		if currentUser.Role != "admin" {
			respWriter.HTTPJson(
				c, 403, "forbidden", "admin only", nil,
			)
			return
		}

		c.Next()
	}
}
