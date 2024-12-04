package validator_util

import (
	"errors"
	"time"
)

func ValidateAuthorBirthdate(birthdate string) error {
	if birthdate == "" {
		return errors.New("birthdate cannot be empty")
	}

	_, err := time.Parse("02-01-2006", birthdate)
	if err != nil {
		return errors.New("invalid birthdate format")
	}

	return nil
}

func ValidateAuthorFirstName(firstName string) error {
	if firstName == "" {
		return errors.New("firstname cannot be empty")
	}

	return nil
}
