package repository

import (
	"book_service/domain/dto"
	"book_service/domain/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BookBorrowRepo struct {
	db *gorm.DB
}

type IBookBorrowRepo interface {
	Create(book *model.BookBorrow) error
	GetByUUID(uuid string) (*model.BookBorrow, error)
	Update(book *model.BookBorrow) error
	Delete(id string) error
	GetList(
		ctx context.Context,
		params dto.BookBorrowRepo_GetListParams,
	) ([]model.BookBorrow, error)
	CountGetList(
		ctx context.Context,
		params dto.BookBorrowRepo_GetListParams,
	) (int64, error)
}

func NewBookBorrowRepo(db *gorm.DB) IBookBorrowRepo {
	return &BookBorrowRepo{
		db: db,
	}
}

func (repo *BookBorrowRepo) Create(bookBorrow *model.BookBorrow) error {
	err := repo.db.Create(bookBorrow).Error
	if err != nil {
		return errors.New("failed to create: " + err.Error())
	}
	return err
}

func (repo *BookBorrowRepo) GetByUUID(uuid string) (*model.BookBorrow, error) {
	var bookBorrow model.BookBorrow
	if err := repo.db.First(&bookBorrow, "uuid = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get: " + err.Error())
	}
	return &bookBorrow, nil
}

func (repo *BookBorrowRepo) Update(bookBorrow *model.BookBorrow) error {
	err := repo.db.Save(bookBorrow).Error
	return err
}

func (repo *BookBorrowRepo) Delete(id string) error {
	err := repo.db.Delete(&model.BookBorrow{}, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("not found")
		}
		return errors.New("failed to delete: " + err.Error())
	}
	return err
}

func (repo *BookBorrowRepo) GetList(
	ctx context.Context,
	params dto.BookBorrowRepo_GetListParams,
) ([]model.BookBorrow, error) {
	// validate param
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		return nil, fmt.Errorf("invalid sort order")
	}

	var models []model.BookBorrow

	tx := repo.db.Model(&models)

	if params.AuthorUUID != "" {
		tx = tx.Where("author_uuid = ?", params.AuthorUUID)
	}

	if params.Query != "" {
		if params.QueryBy != "" {
			tx = tx.Where("? LIKE ?", params.QueryBy, "%"+params.Query+"%")
		} else {
			tx = tx.Where(
				`
					Title LIKE ?
				`,
				"%"+params.Query+"%",
			)
		}
	}

	if params.Page > 0 && params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		tx = tx.Offset(offset).Limit(params.Limit)
	}

	if params.SortOrder != "" && params.SortBy != "" {
		tx = tx.Order(fmt.Sprintf("%s %s", params.SortBy, params.SortOrder))
	}

	err := tx.Find(&models).Error
	if err != nil {
		return nil, errors.New("failed to get: " + err.Error())
	}

	return models, nil
}

func (repo *BookBorrowRepo) CountGetList(
	ctx context.Context,
	params dto.BookBorrowRepo_GetListParams,
) (int64, error) {
	tx := repo.db.Model(&model.BookBorrow{})

	if params.AuthorUUID != "" {
		tx = tx.Where("author_uuid = ?", params.AuthorUUID)
	}

	if params.Query != "" {
		if params.QueryBy != "" {
			tx = tx.Where("? LIKE ?", params.QueryBy, "%"+params.Query+"%")
		} else {
			tx = tx.Where(
				`
					Title LIKE ?
				`,
				"%"+params.Query+"%",
			)
		}
	}

	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return 0, errors.New("failed to count: " + err.Error())
	}

	return count, nil
}
