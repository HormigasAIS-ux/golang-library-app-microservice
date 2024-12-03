package rest

import (
	"auth_service/config"
	"auth_service/domain/dto"
	"auth_service/domain/model"
	"auth_service/interface/rest/handler"
	"auth_service/middleware"
	"auth_service/repository"
	ucase "auth_service/usecase"
	"auth_service/utils/http_response"
	"context"

	_ "auth_service/docs"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger = logging.MustGetLogger("rest")

func SetupServer(router *gin.Engine) {
	_ = context.Background()

	responseWriter := http_response.NewHttpResponseWriter()
	gormDB := config.NewPostgresqlDB()

	// migrations
	err := gormDB.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
	)
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	// repositories
	userRepo := repository.NewUserRepo(gormDB)
	refreshTokenRepo := repository.NewRefreshTokenRepo(gormDB)

	// ucases
	authUcase := ucase.NewAuthUcase(userRepo, refreshTokenRepo)

	// handlers
	authHandler := handler.NewAuthHandler(responseWriter, authUcase)
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

	secureRouter := router.Group("/")
	{
		secureRouter.Use(middleware.AuthMiddleware(responseWriter, authUcase))
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
