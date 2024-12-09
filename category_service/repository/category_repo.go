package repository

import (
	"category_service/domain/dto"
	"category_service/domain/model"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

type ICategoryRepo interface {
	Create(category *model.Category) error
	GetByUUID(uuid string) (*model.Category, error)
	Update(category *model.Category) error
	Delete(id string) error
	GetList(
		params dto.CategoryRepo_GetListParams,
	) ([]model.Category, error)
	CountGetList(
		params dto.CategoryRepo_GetListParams,
	) (int64, error)
}

func NewCategoryRepo(db *gorm.DB) ICategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

func (repo *CategoryRepo) Create(category *model.Category) error {
	err := repo.db.Create(category).Error
	if err != nil {
		return errors.New("failed to create: " + err.Error())
	}
	return err
}

func (repo *CategoryRepo) GetByUUID(uuid string) (*model.Category, error) {
	var category model.Category
	if err := repo.db.First(&category, "uuid = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get: " + err.Error())
	}
	return &category, nil
}

func (repo *CategoryRepo) Update(category *model.Category) error {
	err := repo.db.Save(category).Error
	return err
}

func (repo *CategoryRepo) Delete(id string) error {
	err := repo.db.Delete(&model.Category{}, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("not found")
		}
		return errors.New("failed to delete: " + err.Error())
	}
	return err
}

func (repo *CategoryRepo) GetList(
	params dto.CategoryRepo_GetListParams,
) ([]model.Category, error) {
	// validate param
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		return nil, fmt.Errorf("invalid sort order")
	}

	var models []model.Category

	tx := repo.db.Model(&models)

	if params.Query != "" {
		if params.QueryBy != "" {
			tx = tx.Where("? LIKE ?", params.QueryBy, "%"+params.Query+"%")
		} else {
			var conditions []string
			var args []interface{}
			tmp := model.Category{}
			for _, field := range tmp.GetQueriableFields() {
				conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
				args = append(args, "%"+params.Query+"%")
			}

			tx = tx.Where(strings.Join(conditions, " OR "), args...)
		}
	}

	if params.Offset > 0 && params.Limit > 0 {
		offset := (params.Offset - 1) * params.Limit
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

func (repo *CategoryRepo) CountGetList(
	params dto.CategoryRepo_GetListParams,
) (int64, error) {
	tx := repo.db.Model(&model.Category{})

	if params.Query != "" {
		if params.QueryBy != "" {
			tx = tx.Where("? LIKE ?", params.QueryBy, "%"+params.Query+"%")
		} else {
			var conditions []string
			var args []interface{}
			tmp := model.Category{}
			for _, field := range tmp.GetQueriableFields() {
				conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
				args = append(args, "%"+params.Query+"%")
			}

			tx = tx.Where(strings.Join(conditions, " OR "), args...)
		}
	}

	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return 0, errors.New("failed to count: " + err.Error())
	}

	return count, nil
}
