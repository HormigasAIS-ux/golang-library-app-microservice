package repository

import (
	"auth_service/domain/model"
	"errors"

	"gorm.io/gorm"
)

type RefreshTokenRepo struct {
	db *gorm.DB
}

type IRefreshTokenRepo interface {
	Create(refresh_token *model.RefreshToken) error
	GetByToken(token string) (*model.RefreshToken, error)
	Update(refresh_token *model.RefreshToken) error
	Delete(id string) error
	InvalidateManyByUserUUID(userUUID string) error
}

func NewRefreshTokenRepo(db *gorm.DB) IRefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (repo *RefreshTokenRepo) Create(refresh_token *model.RefreshToken) error {
	err := repo.db.Create(refresh_token).Error
	return err
}

func (repo *RefreshTokenRepo) GetByToken(token string) (*model.RefreshToken, error) {
	var refresh_token model.RefreshToken
	if err := repo.db.First(&refresh_token, "token = ?", token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("refresh_token not found")
		}
		return nil, err
	}
	return &refresh_token, nil
}

func (repo *RefreshTokenRepo) GetByrefresh_tokenname(refresh_tokenname string) (*model.RefreshToken, error) {
	var refresh_token model.RefreshToken
	if err := repo.db.First(&refresh_token, "refresh_tokenname = ?", refresh_tokenname).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("refresh_token not found")
		}
		return nil, err
	}
	return &refresh_token, nil
}

func (repo *RefreshTokenRepo) GetByEmail(email string) (*model.RefreshToken, error) {
	var refresh_token model.RefreshToken
	if err := repo.db.First(&refresh_token, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("refresh_token not found")
		}
		return nil, err
	}
	return &refresh_token, nil
}

func (repo *RefreshTokenRepo) Update(refresh_token *model.RefreshToken) error {
	err := repo.db.Save(refresh_token).Error
	return err
}

func (repo *RefreshTokenRepo) Delete(id string) error {
	err := repo.db.Delete(&model.RefreshToken{}, "id = ?", id).Error
	return err
}

func (repo *RefreshTokenRepo) InvalidateManyByUserUUID(userUUID string) error {
	err := repo.db.Model(&model.RefreshToken{}).Where("invalid = ? AND user_uuid = ?", false, userUUID).Update("invalid", true).Error
	return err
}
