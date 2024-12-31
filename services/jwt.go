package services

import "konbini/utils"

const (
	// JWT_TOKEN_KEY_SIZE represents the size of the byte array for each JWT key/salt
	JWT_TOKEN_KEY_SIZE int = 32
)

// GetRandomJWTKey generates a cryptographically secure byte array of size 32.
func GetRandomJWTKey() ([]byte, error) {
	return utils.RandomBytes(JWT_TOKEN_KEY_SIZE)
}
