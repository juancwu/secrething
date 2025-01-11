package utils

import "strings"

// IsUniqueViolationErr checks if the error is a unique constraint violation
func IsUniqueViolationErr(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
