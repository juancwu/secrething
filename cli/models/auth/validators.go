package auth

import (
	"errors"
	"strings"
)

func validateEmail(s string) error {
	if s == "" {
		return errors.New("Email cannot be empty")
	}

	// shadow validation
	if !strings.Contains(s, "@") {
		return errors.New("Invalid email")
	}
	return nil
}

func validateNickname(s string) error {
	if len(s) < 3 {
		return errors.New("Nickname must be at least 3 characters long")
	}
	return nil
}

func validatePasswords(s ...string) error {
	if len(s) == 1 {
		if len(s[0]) < 12 {
			return errors.New("Password must be at least 12 characters long")
		}
		return nil
	}

	if s[0] != s[1] {
		return errors.New("Passwords do not match")
	}

	return nil
}
