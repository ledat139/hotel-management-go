package validator

import (
	"errors"
	"strings"
)

func ValidateCreateStaffInput(fullName, phone, password string) error {
	if strings.TrimSpace(fullName) == "" || len(fullName) < 2 {
		return errors.New("error.invalid_full_name")
	}
	if len(phone) < 10 || len(phone) > 11 {
		return errors.New("error.invalid_phone_number")
	}
	if len(password) < 6 {
		return errors.New("error.password_must_be_at_least_6_characters")
	}
	return nil
}
