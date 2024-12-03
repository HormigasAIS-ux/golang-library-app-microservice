package model

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	Token     string     `gorm:"unique;not null" json:"token"`
	UserUUID  string     `gorm:"not null" json:"user_uuid"`
	UsedAt    *time.Time `json:"used_at"`
	ExpiredAt *time.Time `json:"expired_at"`
	Invalid   bool       `json:"invalid"`

	User User `gorm:"foreignKey:UserUUID;constraint:OnDelete:CASCADE;"`
}
