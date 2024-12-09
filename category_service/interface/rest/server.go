package rest

import (
	"category_service/config"
	"category_service/domain/dto"
	interface_pkg "category_service/interface"
	rest_handler "category_service/interface/rest/handler"
	rest_middleware "category_service/interface/rest/middleware"
	"category_service/utils/http_response"
	"fmt"

	_ "category_service/docs"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger = logging.MustGetLogger("main")

func SetupServer(commonDependencies interface_pkg.CommonDependency) {
	router := gin.Default()

	respWriter := http_response.NewHttpResponseWriter()

	// handlers
	categoryHandler := rest_handler.NewCategoryHandler(
		commonDependencies.CategoryUcase,
		respWriter,
	)

	// middlewares
	authMiddleware := rest_middleware.AuthMiddleware(respWriter)
	authMiddlewareAdminOnly := rest_middleware.AuthAdminOnlyMiddleware(respWriter)

	// register routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, dto.BaseJSONResp{
			Code:    200,
			Message: "pong",
		})
	})

	secureRouter := router.Group("")
	secureRouter.Use(authMiddleware)
	// secured
	{
		// /books
		bookRouter := secureRouter.Group("/books")
		{
			bookRouter.PATCH("/:category_uuid", categoryHandler.PatchCategory)   // owner only
			bookRouter.DELETE("/:category_uuid", categoryHandler.DeleteCategory) // owner only
			bookRouter.GET("/:category_uuid", categoryHandler.GetCategoryDetail)
			bookRouter.GET("", categoryHandler.GetCategoryList)

			bookRouterAdminOnly := bookRouter.Group("", authMiddlewareAdminOnly)
			{
				bookRouterAdminOnly.POST("", categoryHandler.CreateBook)
			}
		}
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(fmt.Sprintf("%s:%d", config.Envs.HOST, config.Envs.PORT))
}
