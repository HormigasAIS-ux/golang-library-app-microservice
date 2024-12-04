package repository

import (
	"author_service/domain/dto"
	"author_service/domain/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AuthorRepo struct {
	db *gorm.DB
}

type IAuthorRepo interface {
	Create(author *model.Author) error
	GetByUUID(uuid string) (*model.Author, error)
	Update(author *model.Author) error
	Delete(id string) error
	GetList(
		ctx context.Context,
		params dto.AuthorRepo_GetListParams,
	) ([]model.Author, error)
	CountGetList(
		ctx context.Context,
		params dto.AuthorRepo_GetListParams,
	) (int64, error)
}

func NewAuthorRepo(db *gorm.DB) IAuthorRepo {
	return &AuthorRepo{
		db: db,
	}
}

func (repo *AuthorRepo) Create(author *model.Author) error {
	err := repo.db.Create(author).Error
	if err != nil {
		return errors.New("failed to create author")
	}
	return err
}

func (repo *AuthorRepo) GetByUUID(uuid string) (*model.Author, error) {
	var author model.Author
	if err := repo.db.First(&author, "uuid = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("author not found")
		}
		return nil, errors.New("failed to get author")
	}
	return &author, nil
}

func (repo *AuthorRepo) Update(author *model.Author) error {
	err := repo.db.Save(author).Error
	return err
}

func (repo *AuthorRepo) Delete(id string) error {
	err := repo.db.Delete(&model.Author{}, "id = ?", id).Error
	return err
}

func (repo *AuthorRepo) GetList(
	ctx context.Context,
	params dto.AuthorRepo_GetListParams,
) ([]model.Author, error) {
	// validate param
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		return nil, fmt.Errorf("invalid sort order")
	}

	var models []model.Author

	tx := repo.db.Model(&models)

	if params.Query != "" {
		if params.QueryBy != "" {
			tx = tx.Where("? LIKE ?", params.QueryBy, "%"+params.Query+"%")
		} else {
			tx = tx.Where(
				`
					FirstName LIKE ?
					OR LastName LIKE ?
					OR BirthDate LIKE ?
					OR CONCAT(FirstName, ' ', FirstName) LIKE ?
				`,
				"%"+params.Query+"%",
				"%"+params.Query+"%",
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

func (repo *AuthorRepo) CountGetList(
	ctx context.Context,
	params dto.AuthorRepo_GetListParams,
) (int64, error) {
	tx := repo.db.Model(&model.Author{})

	if params.Query != "" {
		if params.QueryBy != "" {
			tx = tx.Where("? LIKE ?", params.QueryBy, "%"+params.Query+"%")
		} else {
			tx = tx.Where(
				`
					FirstName LIKE ?
					OR LastName LIKE ?
					OR BirthDate LIKE ?
					OR CONCAT(FirstName, ' ', FirstName) LIKE ?
				`,
				"%"+params.Query+"%",
				"%"+params.Query+"%",
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
