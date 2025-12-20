package validation

import (
	"errors"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

func ValidateName(name string, fieldName string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New(fieldName + " is required")
	}
	if len(name) < 2 {
		return errors.New(fieldName + " must be at least 2 characters long")
	}
	return nil
}

func ValidatePhoneNumber(phone string) error {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return errors.New("phone number is required")
	}
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone number format")
	}
	return nil
}
