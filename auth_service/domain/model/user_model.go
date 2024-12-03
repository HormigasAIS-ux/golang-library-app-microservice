package model

import (
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

	RefreshTokens []RefreshToken `gorm:"foreignKey:UserUUID;references:UUID;" json:"-"`
}
