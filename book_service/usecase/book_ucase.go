package ucase

import (
	"book_service/domain/dto"
	"book_service/domain/model"
	author_pb "book_service/interface/grpc/genproto/author"
	"book_service/repository"
	error_utils "book_service/utils/error"
	"context"

	"github.com/google/uuid"
	"github.com/op/go-logging"
	"google.golang.org/grpc/codes"
)

var logger = logging.MustGetLogger("main")

type BookUcase struct {
	bookRepo   repository.IBookRepo
	authorRepo repository.IAuthorRepo
}

type IBookUcase interface {
	Create(ctx context.Context, currentUser dto.CurrentUser, payload dto.CreateBookReq) (*dto.CreateBookResp, error)
}

func NewBookUcase(
	bookRepo repository.IBookRepo,
	authorRepo repository.IAuthorRepo,
) IBookUcase {
	return &BookUcase{
		bookRepo:   bookRepo,
		authorRepo: authorRepo,
	}
}

func (ucase *BookUcase) Create(ctx context.Context, currentUser dto.CurrentUser, payload dto.CreateBookReq) (*dto.CreateBookResp, error) {
	// get author by user uuid through author service
	getAuthorResp, grpcCode, err := ucase.authorRepo.RpcCreateAuthor(
		ctx, &author_pb.CreateAuthorReq{
			UserUuid: currentUser.UUID,
		},
	)

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
	books, err := ucase.bookRepo.GetList(
		ctx, dto.BookRepo_GetListParams{
			AuthorUUID: getAuthorResp.Uuid,
			Query:      payload.Title,
			QueryBy:    "title",
			Limit:      1,
			Page:       1,
		},
	)

	if len(books) > 0 {
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
