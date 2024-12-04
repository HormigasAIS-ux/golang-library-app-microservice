package model

import (
	validator_util "author_service/utils/validator/author"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	UUID      uuid.UUID `gorm:"type:uuid;unique;not null" json:"uuid"`
	UserUUID  uuid.UUID `gorm:"type:uuid;not null" json:"user_uuid"`
	FirstName string    `gorm:"type:text;not null" json:"first_name"`
	LastName  string    `gorm:"type:text" json:"last_name"`
	BirthDate *string   `gorm:"type:text" json:"birth_date"`
	Bio       *string   `gorm:"type:text" json:"bio"`
}

func (u *Author) Validate() (err error) {
	// birthdate
	if u.BirthDate != nil {
		err = validator_util.ValidateAuthorBirthdate(*u.BirthDate)
		if err != nil {
			return errors.New("author validation error: " + err.Error())
		}
	}

	// firstname
	err = validator_util.ValidateAuthorFirstName(u.FirstName)
	if err != nil {
		return errors.New("author validation error: " + err.Error())
	}

	return nil
}
