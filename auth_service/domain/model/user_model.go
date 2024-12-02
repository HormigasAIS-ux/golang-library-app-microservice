package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     string `gorm:"unique;not null" json:"uuid"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`

	RefreshTokens []RefreshToken `gorm:"foreignKey:UserUUID" json:"-"`
}
