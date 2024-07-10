package util

import "github.com/google/uuid"

// IsValidUUIDv4 is a utility function to verify if a given string is a valid UUID v4.
func IsValidUUIDv4(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
