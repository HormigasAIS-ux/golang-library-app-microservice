package rest

import (
	"author_service/config"
	"author_service/domain/dto"
	interface_pkg "author_service/interface"
	rest_handler "author_service/interface/rest/handler"
	rest_middleware "author_service/interface/rest/middleware"
	"author_service/utils/http_response"
	"fmt"

	_ "author_service/docs"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger = logging.MustGetLogger("rest")

func SetupServer(commonDependencies interface_pkg.CommonDependency) {
	router := gin.Default()

	respWriter := http_response.NewHttpResponseWriter()

	// handlers
	authorHandler := rest_handler.NewAuthorHandler(
		commonDependencies.AuthorUcase,
		respWriter,
	)

	// register routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, dto.BaseJSONResp{
			Code:    200,
			Message: "pong",
		})
	})

	secureRouter := router.Group("/author")
	secureRouter.Use(rest_middleware.AuthMiddleware(respWriter, commonDependencies.AuthorRepo))
	{
		secureRouter.GET("", authorHandler.GetList)
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(fmt.Sprintf("%s:%d", config.Envs.HOST, config.Envs.PORT))
}
