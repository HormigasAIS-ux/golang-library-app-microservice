package rest

import (
	"auth_service/config"
	"auth_service/domain/dto"
	interface_pkg "auth_service/interface"
	"auth_service/interface/rest/handler"
	"auth_service/utils/http_response"
	"fmt"

	_ "auth_service/docs"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger = logging.MustGetLogger("rest")

func SetupServer(commonDependencies interface_pkg.CommonDependency) {
	router := gin.Default()

	responseWriter := http_response.NewHttpResponseWriter()

	// handlers
	authHandler := handler.NewAuthHandler(responseWriter, commonDependencies.AuthUcase)
	_ = authHandler

	// register routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, dto.BaseJSONResp{
			Code:    200,
			Message: "pong",
		})
	})
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/check-token", authHandler.CheckToken)
	router.POST("/auth/refresh-token", authHandler.RefreshToken)

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(fmt.Sprintf("%s:%d", config.Envs.HOST, config.Envs.PORT))
}
