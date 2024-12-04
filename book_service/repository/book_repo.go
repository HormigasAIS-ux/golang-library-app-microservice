package repository

import (
	"book_service/domain/dto"
	"book_service/domain/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BookRepo struct {
	db *gorm.DB
}

type IBookRepo interface {
	Create(book *model.Book) error
	GetByUUID(uuid string) (*model.Book, error)
	Update(book *model.Book) error
	Delete(id string) error
	GetList(
		ctx context.Context,
		params dto.BookRepo_GetListParams,
	) ([]model.Book, error)
	CountGetList(
		ctx context.Context,
		params dto.BookRepo_GetListParams,
	) (int64, error)
}

func NewBookRepo(db *gorm.DB) IBookRepo {
	return &BookRepo{
		db: db,
	}
}

func (repo *BookRepo) Create(book *model.Book) error {
	err := repo.db.Create(book).Error
	if err != nil {
		return errors.New("failed to create book")
	}
	return err
}

func (repo *BookRepo) GetByUUID(uuid string) (*model.Book, error) {
	var book model.Book
	if err := repo.db.First(&book, "uuid = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get book")
	}
	return &book, nil
}

func (repo *BookRepo) Update(book *model.Book) error {
	err := repo.db.Save(book).Error
	return err
}

func (repo *BookRepo) Delete(id string) error {
	err := repo.db.Delete(&model.Book{}, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("not found")
		}
		return errors.New("failed to delete book")
	}
	return err
}

func (repo *BookRepo) GetList(
	ctx context.Context,
	params dto.BookRepo_GetListParams,
) ([]model.Book, error) {
	// validate param
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		return nil, fmt.Errorf("invalid sort order")
	}

	var models []model.Book

	tx := repo.db.Model(&models)

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
		return nil, err
	}

	return models, nil
}

func (repo *BookRepo) CountGetList(
	ctx context.Context,
	params dto.BookRepo_GetListParams,
) (int64, error) {
	tx := repo.db.Model(&model.Book{})

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
		return 0, err
	}

	return count, nil
}
