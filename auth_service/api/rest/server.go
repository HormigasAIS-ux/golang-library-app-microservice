package rest

import (
	"auth_service/api/rest/handler"
	"auth_service/config"
	"auth_service/domain/dto"
	"auth_service/middleware"
	"auth_service/repository"
	ucase "auth_service/usecase"
	"auth_service/utils/http_response"
	"context"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupServer(router *gin.Engine) {
	_ = context.Background()

	responseWriter := http_response.NewResponseWriter()
	gormDB := config.NewPostgresqlDB()

	// repositories
	userRepo := repository.NewUserRepo(gormDB)
	refreshTokenRepo := repository.NewRefreshTokenRepo(gormDB)

	// ucases
	authUcase := ucase.NewAuthUcase(userRepo, refreshTokenRepo)

	// handlers
	authHandler := handler.NewAuthHandler(responseWriter, authUcase)
	_ = authHandler
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, dto.BaseJSONResp{
			Code:    200,
			Message: "pong",
		})
	})

	secureRouter := router.Group("/")
	{
		secureRouter.Use(middleware.AuthMiddleware(responseWriter, authUcase))
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
