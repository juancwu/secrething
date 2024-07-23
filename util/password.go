package util

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// customValidator represents the custom validator the echo server uses to validate data.
type customValidator struct {
	validator *validator.Validate
}

// Validate method to satisfy the echo validator interface.
func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// ValidatePassword is a custom validator to validate passwords.
// It checks that a password is at least 12 characters long,
// it contains at least one special character,
// at least one uppercase letter,
// at least one lowercase letter,
// and at least one digit.
func ValidatePassword(password string) bool {
	n := len(password)
	if n < 12 {
		return false
	}

	// check for special character
	if matched, _ := regexp.MatchString(`[!@#$%^&*(),.?":{}|<>]`, password); !matched {
		return false
	}

	// check for uppercase letter
	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return false
	}

	// check for lowercase letter
	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return false
	}

	// check for number
	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return false
	}

	return true
}
