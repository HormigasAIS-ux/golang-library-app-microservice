package rest

import (
	"book_service/config"
	"book_service/domain/dto"
	interface_pkg "book_service/interface"
	rest_handler "book_service/interface/rest/handler"
	rest_middleware "book_service/interface/rest/middleware"
	"book_service/utils/http_response"
	"fmt"

	_ "book_service/docs"

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
	bookHandler := rest_handler.NewBookHandler(
		commonDependencies.BookUcase,
		respWriter,
	)

	// middlewares
	authMiddleware := rest_middleware.AuthMiddleware(respWriter)
	// authMiddlewareAdminOnly := rest_middleware.AuthAdminOnlyMiddleware(respWriter)

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
			bookRouter.POST("", bookHandler.Create)
			bookRouter.PATCH("/:book_uuid", bookHandler.PatchBook)
			bookRouter.DELETE("/:book_uuid", bookHandler.DeleteBook)
		}
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/swagger/index.html")
	})

	router.Run(fmt.Sprintf("%s:%d", config.Envs.HOST, config.Envs.PORT))
}
