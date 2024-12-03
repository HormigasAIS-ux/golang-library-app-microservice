package model

import (
	validator_util "auth_service/utils/validator/user"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     uuid.UUID `gorm:"type:uuid;unique;not null" json:"uuid"`
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `gorm:"not null" json:"password"`
	Fullname string    `json:"fullname"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`

	RefreshTokens []RefreshToken `gorm:"foreignKey:UserUUID;references:UUID;" json:"-"`
}

func (u *User) Validate() (err error) {
	// username
	err = validator_util.ValidateUsername(u.Username)
	if err != nil {
		return errors.New("user validation error: " + err.Error())
	}

	// email
	err = validator_util.ValidateEmail(u.Email)
	if err != nil {
		return errors.New("user validation error: " + err.Error())
	}

	// password
	err = validator_util.ValidatePassword(u.Password)
	if err != nil {
		return errors.New("user validation error: " + err.Error())
	}

	return nil
}
