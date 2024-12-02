package repository

import (
	"auth_service/domain/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

type IUserRepo interface {
	Create(user *model.User) error
	GetByUUID(uuid string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id string) error
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) Create(user *model.User) error {
	err := repo.db.Create(user).Error
	return err
}

func (repo *UserRepo) GetByUUID(uuid string) (*model.User, error) {
	var user model.User
	if err := repo.db.First(&user, "uuid = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := repo.db.First(&user, "username = ?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := repo.db.First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) Update(user *model.User) error {
	err := repo.db.Save(user).Error
	return err
}

func (repo *UserRepo) Delete(id string) error {
	err := repo.db.Delete(&model.User{}, "id = ?", id).Error
	return err
}
