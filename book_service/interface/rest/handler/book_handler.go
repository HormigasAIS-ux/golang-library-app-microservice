package rest_handler

import (
	"book_service/domain/dto"
	ucase "book_service/usecase"
	"book_service/utils/helper"
	"book_service/utils/http_response"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

type BookHandler struct {
	bookUcase  ucase.IBookUcase
	respWriter http_response.IHttpResponseWriter
}

type IBookHandler interface {
	Create(ctx *gin.Context)
	PatchBook(ctx *gin.Context)
}

func NewBookHandler(
	bookUcase ucase.IBookUcase,
	respWriter http_response.IHttpResponseWriter,
) IBookHandler {
	return &BookHandler{
		bookUcase:  bookUcase,
		respWriter: respWriter,
	}
}

// @Summary Create new book
// @Router /books [post]
// @Tags Books
// @Param payload body dto.CreateBookReq true "payload"
// @Success 200 {object} dto.BaseJSONResp{data=dto.CreateBookResp}
// @Security BearerAuth
func (handler *BookHandler) Create(ctx *gin.Context) {
	var payload dto.CreateBookReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("invalid payload: %v", err)
		handler.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	currentUser, err := helper.GetCurrentUserFromGinCtx(ctx)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	resp, err := handler.bookUcase.Create(ctx, *currentUser, payload)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	handler.respWriter.HTTPJsonOK(ctx, resp)
}

// @Summary patch book
// @Router /books/{book_uuid} [patch]
// @Tags Books
// @Param payload body dto.PatchBookReq true "payload"
// @Success 200 {object} dto.BaseJSONResp{data=dto.PatchBookRespData}
// @Security BearerAuth
func (handler *BookHandler) PatchBook(ctx *gin.Context) {
	bookUUID := ctx.Param("book_uuid")

	var payload dto.PatchBookReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("invalid payload: %v", err)
		handler.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	currentUser, err := helper.GetCurrentUserFromGinCtx(ctx)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	data, err := handler.bookUcase.PatchBook(ctx, *currentUser, bookUUID, payload)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	handler.respWriter.HTTPJsonOK(ctx, data)
}
