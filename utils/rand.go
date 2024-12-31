package utils

import "crypto/rand"

// RandomBytes generates a cryptographically secure random byte array of given size.
func RandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
