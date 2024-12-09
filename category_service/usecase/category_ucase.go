package ucase

import (
	"category_service/domain/dto"
	"category_service/domain/model"
	book_grpc "category_service/interface/grpc/genproto/book"
	"category_service/repository"
	error_utils "category_service/utils/error"
	"context"

	"github.com/google/uuid"
	"github.com/op/go-logging"
	"google.golang.org/grpc/codes"
)

var logger = logging.MustGetLogger("main")

type CategoryUcase struct {
	categoryRepo          repository.ICategoryRepo
	bookGrpcServiceClient book_grpc.BookServiceClient
}

type ICategoryUcase interface {
	Create(ctx context.Context, currentUser dto.CurrentUser, payload dto.CreateCategoryReq) (*dto.CreateCategoryRespData, error)
	PatchCategory(
		ctx context.Context,
		currentUser dto.CurrentUser,
		categoryUUID string,
		payload dto.PatchCategoryReq,
	) (*dto.PatchCategoryRespData, error)
	DeleteCategory(
		ctx context.Context,
		currentUser dto.CurrentUser,
		categoryUUID string,
	) (*dto.DeleteCategoryRespData, error)
	GetCategoryDetail(
		ctx context.Context,
		categoryUUID string,
	) (*dto.GetCategoryDetailRespData, error)
	GetListCategory(
		ctx context.Context,
		params dto.GetCategoryListReq,
	) (*dto.GetListCategoryRespData, error)
}

func NewCategoryUcase(
	categoryRepo repository.ICategoryRepo,
	bookGrpcServiceClient book_grpc.BookServiceClient,
) ICategoryUcase {
	return &CategoryUcase{
		categoryRepo:          categoryRepo,
		bookGrpcServiceClient: bookGrpcServiceClient,
	}
}

func (ucase *CategoryUcase) Create(ctx context.Context, currentUser dto.CurrentUser, payload dto.CreateCategoryReq) (*dto.CreateCategoryRespData, error) {
	// validate input
	if payload.Name == "" {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.InvalidArgument,
			Message:  "invalid input",
			Detail:   "name cannot be empty",
		}
	}

	// check title exists by author uuid
	if categories, _ := ucase.categoryRepo.GetList(
		dto.CategoryRepo_GetListParams{
			Query:   payload.Name,
			QueryBy: "name",
			Limit:   1,
			Offset:  0,
		},
	); len(categories) > 0 {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.AlreadyExists,
			Message:  "conflict",
			Detail:   "category already exists",
		}
	}

	// create category
	parsedUserUUID, _ := uuid.Parse(currentUser.UUID)
	newCategory := &model.Category{
		CreatedBy: parsedUserUUID,
		Name:      payload.Name,
	}

	err := ucase.categoryRepo.Create(newCategory)
	if err != nil {
		logger.Debugf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	return &dto.CreateCategoryRespData{
		UUID:      newCategory.UUID.String(),
		CreatedBy: newCategory.CreatedAt.String(),
		Name:      newCategory.Name,
		CreatedAt: newCategory.CreatedAt,
		UpdatedAt: newCategory.UpdatedAt,
	}, nil
}

func (ucase *CategoryUcase) PatchCategory(
	ctx context.Context,
	currentUser dto.CurrentUser,
	categoryUUID string,
	payload dto.PatchCategoryReq,
) (*dto.PatchCategoryRespData, error) {
	// find category
	category, err := ucase.categoryRepo.GetByUUID(categoryUUID)
	if err != nil {
		if err.Error() == "not found" {
			logger.Errorf("err: %v", err)
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				GrpcCode: codes.NotFound,
				Message:  "not found",
				Detail:   err,
			}
		} else {
			logger.Errorf("err: %v", err)
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				GrpcCode: codes.Internal,
				Message:  "internal server error",
				Detail:   err,
			}
		}
	}

	// validate user
	if category.CreatedBy.String() != currentUser.UUID {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 403,
			GrpcCode: codes.PermissionDenied,
			Message:  "forbidden",
			Detail:   "forbidden",
		}
	}

	if payload.Name != nil {
		category.Name = *payload.Name
	}

	// update category
	err = ucase.categoryRepo.Update(category)
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err.Error(),
		}
	}

	return &dto.PatchCategoryRespData{
		UUID:      category.UUID.String(),
		CreatedBy: category.CreatedAt.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (ucase *CategoryUcase) DeleteCategory(
	ctx context.Context,
	currentUser dto.CurrentUser,
	categoryUUID string,
) (*dto.DeleteCategoryRespData, error) {
	// find category
	category, err := ucase.categoryRepo.GetByUUID(categoryUUID)
	if err != nil {
		if err.Error() == "not found" {
			logger.Errorf("err: %v", err)
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				GrpcCode: codes.NotFound,
				Message:  "not found",
				Detail:   err,
			}
		} else {
			logger.Errorf("err: %v", err)
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				GrpcCode: codes.Internal,
				Message:  "internal server error",
				Detail:   err,
			}
		}
	}

	// validate user
	if category.CreatedBy.String() != currentUser.UUID {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 403,
			GrpcCode: codes.PermissionDenied,
			Message:  "forbidden",
			Detail:   "forbidden",
		}
	}

	// delete category
	err = ucase.categoryRepo.Delete(categoryUUID)
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	return &dto.DeleteCategoryRespData{
		UUID:      category.UUID.String(),
		CreatedBy: category.CreatedAt.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (ucase *CategoryUcase) GetCategoryDetail(
	ctx context.Context,
	categoryUUID string,
) (*dto.GetCategoryDetailRespData, error) {
	// find category
	category, err := ucase.categoryRepo.GetByUUID(categoryUUID)
	if err != nil {
		if err.Error() == "not found" {
			logger.Errorf("err: %v", err)
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				GrpcCode: codes.NotFound,
				Message:  "not found",
				Detail:   err,
			}
		} else {
			logger.Errorf("err: %v", err)
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				GrpcCode: codes.Internal,
				Message:  "internal server error",
				Detail:   err,
			}
		}
	}

	// TODO: get book total by category uuid through book service

	return &dto.GetCategoryDetailRespData{
		UUID:      category.UUID.String(),
		CreatedBy: category.CreatedAt.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (ucase *CategoryUcase) GetListCategory(
	ctx context.Context,
	params dto.GetCategoryListReq,
) (*dto.GetListCategoryRespData, error) {
	// prepare queryBy
	queryBy := params.QueryBy
	if queryBy == "any" {
		queryBy = ""
	}

	// prepare offset
	var offset int
	if params.Limit == 0 {
		offset = 0
	} else {
		offset = (params.Page - 1) * params.Limit
	}

	// get list
	categories, err := ucase.categoryRepo.GetList(
		dto.CategoryRepo_GetListParams{
			Query:     params.Query,
			QueryBy:   queryBy,
			Offset:    offset,
			Limit:     params.Limit,
			SortOrder: params.SortOrder,
			SortBy:    params.SortBy,
		},
	)

	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	// count
	count, err := ucase.categoryRepo.CountGetList(
		dto.CategoryRepo_GetListParams{
			Query:     params.Query,
			QueryBy:   queryBy,
			Offset:    offset,
			Limit:     params.Limit,
			SortOrder: params.SortOrder,
			SortBy:    params.SortBy,
		},
	)

	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	// TODO: get bulk book total by category uuids through book service

	res := &dto.GetListCategoryRespData{}
	res.Set(params.Page, params.Limit, count)
	for _, category := range categories {
		res.Data = append(res.Data, dto.GetListCategoryRespDataItem{
			UUID:      category.UUID.String(),
			CreatedBy: category.CreatedAt.String(),
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		})
	}

	return res, nil
}
