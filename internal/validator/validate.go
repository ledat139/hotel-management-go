package validator

import (
	"errors"
	"strings"
)

func ValidateCreateStaffInput(fullName, phone string) error {
	if strings.TrimSpace(fullName) == "" || len(fullName) < 2 {
		return errors.New("error.invalid_full_name")
	}
	if len(phone) < 10 || len(phone) > 11 {
		return errors.New("error.invalid_phone_number")
	}
	return nil
}
