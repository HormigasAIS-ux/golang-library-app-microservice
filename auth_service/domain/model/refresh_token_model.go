package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	Token     string     `gorm:"type:text;unique;not null" json:"token"`
	UserUUID  uuid.UUID  `gorm:"type:uuid;not null" json:"user_uuid"`
	UsedAt    *time.Time `json:"used_at"`
	ExpiredAt *time.Time `json:"expired_at"`
	Invalid   bool       `json:"invalid"`

	User User `gorm:"foreignKey:UserUUID;references:UUID;constraint:OnDelete:CASCADE;"`
}
