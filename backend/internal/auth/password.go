package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CompareHashes(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
