package ucase

import (
	"book_service/domain/dto"
	"book_service/repository"
	"context"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

type BookUcase struct {
	bookRepo repository.IBookRepo
}

type IBookUcase interface {
}

func NewBookUcase(
	bookRepo repository.IBookRepo,
) IBookUcase {
	return &BookUcase{
		bookRepo: bookRepo,
	}
}

func (repo *BookUcase) Create(ctx context.Context, currentUser dto.CurrentUser, payload dto.CreateBookReq) error {
	// check title exists by author uuid
}
