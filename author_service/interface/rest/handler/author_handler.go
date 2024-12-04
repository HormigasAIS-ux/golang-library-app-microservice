package rest_handler

import (
	"author_service/domain/dto"
	ucase "author_service/usecase"
	"author_service/utils/http_response"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	authorUcase ucase.IAuthorUcase
	respWriter  http_response.IHttpResponseWriter
}

type IAuthorHandler interface {
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
