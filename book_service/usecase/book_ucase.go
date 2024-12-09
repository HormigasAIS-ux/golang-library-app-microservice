package ucase

import (
	"book_service/domain/dto"
	"book_service/domain/model"
	author_grpc "book_service/interface/grpc/genproto/author"
	"book_service/repository"
	error_utils "book_service/utils/error"
	"context"

	"github.com/google/uuid"
	"github.com/op/go-logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = logging.MustGetLogger("main")

type BookUcase struct {
	bookRepo                repository.IBookRepo
	authorGrpcServiceClient author_grpc.AuthorServiceClient
}

type IBookUcase interface {
	Create(ctx context.Context, currentUser dto.CurrentUser, payload dto.CreateBookReq) (*dto.CreateBookResp, error)
	PatchBook(
		ctx context.Context,
		currentUser dto.CurrentUser,
		bookUUID string,
		payload dto.PatchBookReq,
	) (*dto.PatchBookRespData, error)
	DeleteBook(
		ctx context.Context,
		currentUser dto.CurrentUser,
		bookUUID string,
	) (*dto.DeleteBookRespData, error)
	GetBookTotalByAuthorUUID(ctx context.Context, authorUUID string) (int64, error)
	BulkGetBookTotalByAuthorUUIDs(
		ctx context.Context,
		payload dto.BulkGetBookTotalByAuthorUUIDsReq,
	) ([]dto.BulkGetBookTotalByAuthorUUIDsRespDataItem, error)
}

func NewBookUcase(
	bookRepo repository.IBookRepo,
	authorGrpcServiceClient author_grpc.AuthorServiceClient,
) IBookUcase {
	return &BookUcase{
		bookRepo:                bookRepo,
		authorGrpcServiceClient: authorGrpcServiceClient,
	}
}

func (ucase *BookUcase) Create(ctx context.Context, currentUser dto.CurrentUser, payload dto.CreateBookReq) (*dto.CreateBookResp, error) {
	// get author by user uuid through author service
	getAuthorResp, err := ucase.authorGrpcServiceClient.GetAuthorByUserUUID(
		ctx, &author_grpc.GetAuthorByUserUUIDReq{
			UserUuid: currentUser.UUID,
		},
	)

	grpcCode := status.Code(err)

	if grpcCode != codes.OK {
		logger.Debugf("grpcCode: %v;\nerr: %v", grpcCode, err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	if getAuthorResp == nil || getAuthorResp.Uuid == "" {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   "getAuthorResp is nil or getAuthorResp.Uuid is empty",
		}
	}

	// check title exists by author uuid
	if books, _ := ucase.bookRepo.GetList(
		dto.BookRepo_GetListParams{
			AuthorUUID: getAuthorResp.Uuid,
			Query:      payload.Title,
			QueryBy:    "title",
			Limit:      1,
			Page:       1,
		},
	); len(books) > 0 {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.AlreadyExists,
			Message:  "conflict",
			Detail:   "book already exists",
		}
	}

	// TODO: validate category through category service

	// create book
	var parsedCategoryUUID *uuid.UUID = nil
	if payload.CategoryUUID != nil {
		tmp, _ := uuid.Parse(*payload.CategoryUUID)
		parsedCategoryUUID = &tmp

	}
	newBook := &model.Book{
		AuthorUUID:   uuid.MustParse(getAuthorResp.Uuid),
		Title:        payload.Title,
		Stock:        payload.Stock,
		CategoryUUID: parsedCategoryUUID,
	}

	err = ucase.bookRepo.Create(newBook)
	if err != nil {
		logger.Debugf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	return &dto.CreateBookResp{
		UUID:       newBook.UUID.String(),
		AuthorUUID: newBook.AuthorUUID.String(),
		CategoryUUID: func() *string {
			if newBook.CategoryUUID == nil {
				return nil
			}
			tmp := newBook.CategoryUUID.String()
			return &tmp
		}(),
		Title:     newBook.Title,
		Stock:     newBook.Stock,
		CreatedAt: newBook.CreatedAt,
		UpdatedAt: newBook.UpdatedAt,
	}, nil
}

func (ucase *BookUcase) PatchBook(
	ctx context.Context,
	currentUser dto.CurrentUser,
	bookUUID string,
	payload dto.PatchBookReq,
) (*dto.PatchBookRespData, error) {
	// find book
	book, err := ucase.bookRepo.GetByUUID(bookUUID)
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

	// get author by user uuid through author service
	getAuthorResp, err := ucase.authorGrpcServiceClient.GetAuthorByUserUUID(
		ctx, &author_grpc.GetAuthorByUserUUIDReq{
			UserUuid: currentUser.UUID,
		},
	)

	grpcCode := status.Code(err)

	if grpcCode != codes.OK {
		logger.Debugf("grpcCode: %v;\nerr: %v", grpcCode, err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	if getAuthorResp == nil || getAuthorResp.Uuid == "" {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   "getAuthorResp is nil or getAuthorResp.Uuid is empty",
		}
	}

	// validate user
	if book.AuthorUUID.String() != getAuthorResp.Uuid {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 403,
			GrpcCode: codes.PermissionDenied,
			Message:  "forbidden",
			Detail:   "forbidden",
		}
	}

	if payload.CategoryUUID != nil {
		// TODO: validate category through category service

		noValue := "no value"
		if payload.CategoryUUID == &noValue {
			book.CategoryUUID = nil
		} else {
			tmp, _ := uuid.Parse(*payload.CategoryUUID)
			book.CategoryUUID = &tmp
		}
	}

	if payload.Title != nil {
		book.Title = *payload.Title
	}

	if payload.Stock != nil {
		book.Stock = *payload.Stock
	}

	// update book
	err = ucase.bookRepo.Update(book)
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err.Error(),
		}
	}

	return &dto.PatchBookRespData{
		UUID:       book.UUID.String(),
		AuthorUUID: book.AuthorUUID.String(),
		CategoryUUID: func() *string {
			if book.CategoryUUID == nil {
				return nil
			}
			tmp := book.CategoryUUID.String()
			return &tmp
		}(),
		Title:     book.Title,
		Stock:     book.Stock,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}, nil
}

func (ucase *BookUcase) DeleteBook(
	ctx context.Context,
	currentUser dto.CurrentUser,
	bookUUID string,
) (*dto.DeleteBookRespData, error) {
	// find book
	book, err := ucase.bookRepo.GetByUUID(bookUUID)
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

	// get author by user uuid through author service
	getAuthorResp, err := ucase.authorGrpcServiceClient.GetAuthorByUserUUID(
		ctx, &author_grpc.GetAuthorByUserUUIDReq{
			UserUuid: currentUser.UUID,
		},
	)

	grpcCode := status.Code(err)

	if grpcCode != codes.OK {
		logger.Debugf("grpcCode: %v;\nerr: %v", grpcCode, err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	if getAuthorResp == nil || getAuthorResp.Uuid == "" {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   "getAuthorResp is nil or getAuthorResp.Uuid is empty",
		}
	}

	// validate user
	if book.AuthorUUID.String() != getAuthorResp.Uuid {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 403,
			GrpcCode: codes.PermissionDenied,
			Message:  "forbidden",
			Detail:   "forbidden",
		}
	}

	// delete book
	err = ucase.bookRepo.Delete(bookUUID)
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	return &dto.DeleteBookRespData{
		UUID:       book.UUID.String(),
		AuthorUUID: book.AuthorUUID.String(),
		CategoryUUID: func() *string {
			if book.CategoryUUID == nil {
				return nil
			}
			tmp := book.CategoryUUID.String()
			return &tmp
		}(),
		Title:     book.Title,
		Stock:     book.Stock,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}, nil
}

func (ucase *BookUcase) GetBookTotalByAuthorUUID(
	ctx context.Context,
	authorUUID string,
) (int64, error) {
	if authorUUID == "" {
		return 0, &error_utils.CustomErr{
			HttpCode: 400,
			GrpcCode: codes.InvalidArgument,
			Message:  "invalid argument",
			Detail:   "authorUUID is empty",
		}
	}

	count, err := ucase.bookRepo.CountGetList(
		dto.BookRepo_GetListParams{
			AuthorUUID: authorUUID,
		},
	)
	if err != nil {
		logger.Errorf("err: %v", err)
		return 0, &error_utils.CustomErr{
			HttpCode: 500,
			GrpcCode: codes.Internal,
			Message:  "internal server error",
			Detail:   err,
		}
	}

	return count, nil
}

func (ucase *BookUcase) BulkGetBookTotalByAuthorUUIDs(
	ctx context.Context,
	payload dto.BulkGetBookTotalByAuthorUUIDsReq,
) ([]dto.BulkGetBookTotalByAuthorUUIDsRespDataItem, error) {
	var results []dto.BulkGetBookTotalByAuthorUUIDsRespDataItem
	for _, authorUUID := range payload.AuthorUUIDs {
		count, err := ucase.bookRepo.CountGetList(
			dto.BookRepo_GetListParams{
				AuthorUUID: authorUUID,
			},
		)
		if err != nil {
			logger.Warningf("err: %v", err)
		}

		results = append(results, dto.BulkGetBookTotalByAuthorUUIDsRespDataItem{
			AuthorUUID: authorUUID,
			Total:      count,
		})
	}

	return results, nil
}
