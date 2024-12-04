package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	UUID         uuid.UUID  `gorm:"type:uuid;unique;not null" json:"uuid"`
	AuthorUUID   uuid.UUID  `gorm:"type:uuid;not null" json:"author_uuid"`
	CategoryUUID *uuid.UUID `gorm:"type:uuid" json:"category_uuid"`
	Title        string     `gorm:"type:text;not null" json:"title"`
	Stock        int64      `json:"stock"`

	BookBorrows []BookBorrow `gorm:"foreignKey:BookUUID;references:UUID;constraint:OnDelete:CASCADE;" json:"-"`
}

type BookBorrow struct {
	gorm.Model
	UUID       uuid.UUID `gorm:"type:uuid;unique;not null" json:"uuid"`
	BookUUID   uuid.UUID `gorm:"type:uuid;not null" json:"book_uuid"`
	UserUUID   uuid.UUID `gorm:"type:uuid;not null" json:"user_uuid"`
	Book       Book      `gorm:"foreignKey:BookUUID;references:UUID;constraint:OnDelete:CASCADE;" json:"-"`
	BorrowDate *string   `json:"borrow_date"`
	ReturnDate *string   `json:"return_date"`
}
