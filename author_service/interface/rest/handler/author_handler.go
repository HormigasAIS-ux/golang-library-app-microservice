package rest_handler

import (
	"author_service/domain/dto"
	ucase "author_service/usecase"
	"author_service/utils/http_response"
	"context"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	authorUcase ucase.IAuthorUcase
	respWriter  http_response.IHttpResponseWriter
}

type IAuthorHandler interface {
	CreateNewAuthor(ctx *gin.Context)
	EditMe(ctx *gin.Context)
	EditAuthor(ctx *gin.Context)
	DeleteAuthor(ctx *gin.Context)
	GetMe(ctx *gin.Context)
	GetAuthorDetail(ctx *gin.Context)
	GetList(ctx *gin.Context)
}

func NewAuthorHandler(
	authorUcase ucase.IAuthorUcase,
	respWriter http_response.IHttpResponseWriter,
) IAuthorHandler {
	return &AuthorHandler{
		authorUcase: authorUcase,
		respWriter:  respWriter,
	}
}

// @Summary Create new author
// @Router /authors [post]
// @Tags Authors
// @Param payload body dto.CreateNewAuthorReq true "payload"
// @Success 200 {object} dto.BaseJSONResp{data=dto.CreateNewAuthorRespData}
func (h *AuthorHandler) CreateNewAuthor(ctx *gin.Context) {
	var payload dto.CreateNewAuthorReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(
			ctx, 400, "invalid request", err.Error(), nil,
		)
		return
	}

	resp, err := h.authorUcase.CreateNewAuthor(context.TODO(), ctx, payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(
			ctx, err,
		)
		return
	}

	h.respWriter.HTTPJsonOK(
		ctx, resp,
	)
}

// @Summary Edit my author profile
// @Router /authors/me [patch]
// @Tags Authors
// @Param payload body dto.EditAuthorReq true "payload"
// @Success 200 {object} dto.BaseJSONResp{data=dto.EditAuthorRespData}
func (h *AuthorHandler) EditMe(ctx *gin.Context) {
	var payload dto.EditAuthorReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(
			ctx, 400, "invalid request", err.Error(), nil,
		)
		return
	}

	resp, err := h.authorUcase.EditAuthor(ctx, "me", payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(
			ctx, err,
		)
		return
	}

	h.respWriter.HTTPJsonOK(
		ctx, resp,
	)
}

// @Summary Edit author
// @Router /authors/{author_uuid} [patch]
// @Tags Authors
// @Param payload body dto.EditAuthorReq true "payload"
// @Success 200 {object} dto.BaseJSONResp{data=dto.EditAuthorRespData}
func (h *AuthorHandler) EditAuthor(ctx *gin.Context) {
	authorUUID := ctx.Param("author_uuid")

	var payload dto.EditAuthorReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(
			ctx, 400, "invalid request", err.Error(), nil,
		)
		return
	}

	resp, err := h.authorUcase.EditAuthor(ctx, authorUUID, payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(
			ctx, err,
		)
		return
	}

	h.respWriter.HTTPJsonOK(
		ctx, resp,
	)
}

// @Summary Delete author
// @Router /authors/{author_uuid} [delete]
// @Tags Authors
// @Success 200 {object} dto.BaseJSONResp{data=dto.DeleteAuthorRespData}
func (h *AuthorHandler) DeleteAuthor(ctx *gin.Context) {
	authorUUID := ctx.Param("author_uuid")

	resp, err := h.authorUcase.DeleteAuthor(ctx, authorUUID)
	if err != nil {
		h.respWriter.HTTPCustomErr(
			ctx, err,
		)
		return
	}

	h.respWriter.HTTPJsonOK(
		ctx, resp,
	)
}

// @Summary Get my author profile detail
// @Router /authors/me [get]
// @Tags Authors
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetAuthorDetailRespData}
func (h *AuthorHandler) GetMe(ctx *gin.Context) {
	resp, err := h.authorUcase.GetAuthorDetail(ctx, "me")
	if err != nil {
		h.respWriter.HTTPCustomErr(
			ctx, err,
		)
		return
	}

	h.respWriter.HTTPJsonOK(
		ctx, resp,
	)
}

// @Summary Get author detail
// @Router /authors/{author_uuid} [get]
// @Tags Authors
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetAuthorDetailRespData}
func (h *AuthorHandler) GetAuthorDetail(ctx *gin.Context) {
	authorUUID := ctx.Param("author_uuid")

	resp, err := h.authorUcase.GetAuthorDetail(ctx, authorUUID)
	if err != nil {
		h.respWriter.HTTPCustomErr(
			ctx, err,
		)
		return
	}

	h.respWriter.HTTPJsonOK(
		ctx, resp,
	)
}

// @Summary Get Author List
// @Router /authors [get]
// @Tags Authors
// @Param query query dto.GetAuthorListReq false "queries"
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetAuthorListRespData}
func (h *AuthorHandler) GetList(ctx *gin.Context) {
	var query dto.GetAuthorListReq
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		h.respWriter.HTTPJson(
			ctx, 400, "invalid request", err.Error(), nil,
		)
		return
	}

	data, count, err := h.authorUcase.GetList(ctx, query)
	if err != nil {
		h.respWriter.HTTPCustomErr(
			ctx, err,
		)
		return
	}

	resp := dto.GetAuthorListRespData{Data: data}
	resp.Set(query.Page, query.Limit, count)

	h.respWriter.HTTPJsonOK(
		ctx, resp,
	)
}
