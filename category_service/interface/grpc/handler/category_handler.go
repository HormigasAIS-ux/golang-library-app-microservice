package handler

import (
	category_grpc "category_service/interface/grpc/genproto/category"
	ucase "category_service/usecase"

	"github.com/op/go-logging"
)

type CategoryServiceHandler struct {
	category_grpc.UnimplementedCategoryServiceServer
	categoryUcase ucase.ICategoryUcase
}

var logger = logging.MustGetLogger("main")

func NewCategoryServiceHandler(categoryUcase ucase.ICategoryUcase) *CategoryServiceHandler {
	handler := &CategoryServiceHandler{categoryUcase: categoryUcase}
	return handler
}
