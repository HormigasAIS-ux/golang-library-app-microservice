package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	UUID      uuid.UUID `gorm:"type:uuid;unique;not null" json:"uuid"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

func (category *Category) GetQueriableFields() []string {
	return []string{"name"}
}
