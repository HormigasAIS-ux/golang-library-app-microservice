package rest_handler

import (
	"category_service/domain/dto"
	ucase "category_service/usecase"
	"category_service/utils/helper"
	"category_service/utils/http_response"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

type CategoryHandler struct {
	categoryUcase ucase.ICategoryUcase
	respWriter    http_response.IHttpResponseWriter
}

type ICategoryHandler interface {
	CreateBook(ctx *gin.Context)
	PatchCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)
	GetCategoryDetail(ctx *gin.Context)
	GetCategoryList(ctx *gin.Context)
}

func NewCategoryHandler(
	categoryUcase ucase.ICategoryUcase,
	respWriter http_response.IHttpResponseWriter,
) ICategoryHandler {
	return &CategoryHandler{
		categoryUcase: categoryUcase,
		respWriter:    respWriter,
	}
}

// @Summary Create new category
// @Router /category [post]
// @Tags Categories
// @Param payload body dto.CreateCategoryReq true "payload"
// @Success 200 {object} dto.BaseJSONResp{data=dto.CreateCategoryRespData}
// @Security BearerAuth
func (handler *CategoryHandler) CreateBook(ctx *gin.Context) {
	var payload dto.CreateCategoryReq
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

	resp, err := handler.categoryUcase.Create(ctx, *currentUser, payload)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	handler.respWriter.HTTPJsonOK(ctx, resp)
}

// @Summary patch category
// @Router /category/{category_uuid} [patch]
// @Tags Categories
// @Param payload body dto.PatchCategoryReq true "payload"
// @Success 200 {object} dto.BaseJSONResp{data=dto.PatchCategoryRespData}
// @Security BearerAuth
func (handler *CategoryHandler) PatchCategory(ctx *gin.Context) {
	categoryUUID := ctx.Param("category_uuid")

	var payload dto.PatchCategoryReq
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

	data, err := handler.categoryUcase.PatchCategory(ctx, *currentUser, categoryUUID, payload)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	handler.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary Delete category
// @Router /category/{category_uuid} [delete]
// @Tags Categories
// @Success 200 {object} dto.BaseJSONResp{data=dto.DeleteCategoryRespData}
// @Security BearerAuth
func (handler *CategoryHandler) DeleteCategory(
	ctx *gin.Context,
) {
	categoryUUID := ctx.Param("category_uuid")

	currentUser, err := helper.GetCurrentUserFromGinCtx(ctx)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	data, err := handler.categoryUcase.DeleteCategory(
		ctx,
		*currentUser,
		categoryUUID,
	)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	handler.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary Get category detail
// @Router /category/{category_uuid} [get]
// @Tags Categories
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetCategoryDetailRespData}
// @Security BearerAuth
func (handler *CategoryHandler) GetCategoryDetail(
	ctx *gin.Context,
) {
	categoryUUID := ctx.Param("category_uuid")

	resp, err := handler.categoryUcase.GetCategoryDetail(ctx, categoryUUID)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	handler.respWriter.HTTPJsonOK(ctx, resp)
}

// @Summary Get category list
// @Router /categories [get]
// @Tags Categories
// @Param query query dto.GetCategoryListReq true "query"
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetListCategoryRespData}
// @Security BearerAuth
func (handler *CategoryHandler) GetCategoryList(
	ctx *gin.Context,
) {
	var queries dto.GetCategoryListReq
	if err := ctx.ShouldBindQuery(&queries); err != nil {
		logger.Errorf("invalid query: %v", err)
		handler.respWriter.HTTPJson(ctx, 400, "invalid query", err.Error(), nil)
		return
	}

	resp, err := handler.categoryUcase.GetListCategory(ctx, queries)
	if err != nil {
		handler.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	handler.respWriter.HTTPJsonOK(ctx, resp)
}
