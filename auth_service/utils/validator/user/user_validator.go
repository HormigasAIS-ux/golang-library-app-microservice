package validator_util

import (
	"errors"
	"strings"
)

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if strings.Contains(username, " ") {
		return errors.New("username cannot contain spaces")
	}
	return nil
}

func ValidateEmail(email string) error {
	if strings.Contains(email, " ") {
		return errors.New("email cannot contain spaces")
	}
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @")
	}
	return nil
}

func ValidateRawPassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if strings.Contains(password, " ") {
		return errors.New("password cannot contain spaces")
	}
	return nil
}
