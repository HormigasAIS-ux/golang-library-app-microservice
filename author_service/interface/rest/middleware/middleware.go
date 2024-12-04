package rest_middleware

import (
	"author_service/config"
	"author_service/repository"
	"author_service/utils/http_response"
	jwt_util "author_service/utils/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(respWriter http_response.IHttpResponseWriter, authRepo repository.IAuthRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		currentUser, err := jwt_util.ValidateJWT(token, config.Envs.JWT_SECRET_KEY)
		if err != nil {
			respWriter.HTTPJson(
				c, 401, "unauthorized", err.Error(), nil,
			)
			return
		}

		c.Set("currentUser", currentUser)
		c.Next()
	}
}
