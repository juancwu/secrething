package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Memory      uint32 = 65536 // 64MB
	argon2SaltLength  uint32 = 16
	argon2KeyLength   uint32 = 32
	argon2Iterations  uint32 = 3
	argon2Parallelism uint8  = 2
)

var (
	ErrInvalidEncodedHash        error = errors.New("Invalid encoded hash value")
	ErrIncompatibleArgon2Version error = errors.New("Incompatible Argon2 version")
)

type HashParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	KeyLength   uint32
}

type DecodedHash struct {
	Params HashParams
	Salt   []byte
	Hash   []byte
}

// Generates an encoded password hash that can be safely stored.
func GeneratePasswordHash(password string) (string, error) {
	params := HashParams{
		Memory:      argon2Memory,
		Iterations:  argon2Iterations,
		Parallelism: argon2Parallelism,
		KeyLength:   argon2KeyLength,
	}
	hash, salt, err := HashPassword(password, params, nil)
	if err != nil {
		return "", err
	}

	// encode the hash and its parameters
	encodedResult := EncodePasswordHash(hash, salt, params)

	return encodedResult, nil
}

// HashPassword hashes the given password using the Argon2 algorithm.
// Optionally pass a salt byte array or nil to generate a random salt.
// Returns the hashed byte array and the salt byte array
func HashPassword(password string, params HashParams, salt []byte) ([]byte, []byte, error) {
	var err error
	if salt == nil {
		salt, err = RandomBytes(argon2SaltLength)
		if err != nil {
			return nil, nil, err
		}
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	return hash, salt, nil
}

func ComparePasswordAndHash(password string, encodedHash string) (bool, error) {
	decodedHash, err := DecodePasswordHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash, _, err := HashPassword(password, decodedHash.Params, decodedHash.Salt)
	if err != nil {
		return false, err
	}

	result := subtle.ConstantTimeCompare(decodedHash.Hash, otherHash)

	return result == 1, nil
}

// EncodePasswordHash encodes the given hash and salt into a string.
func EncodePasswordHash(hash []byte, salt []byte, params HashParams) string {
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedResult := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedResult
}

// DecodePasswordHash decodes the encoded hash into a DecodedHash struct
// which has all the fields needed to further compare passwords.
func DecodePasswordHash(encodedHash string) (*DecodedHash, error) {
	var err error
	var result DecodedHash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, ErrInvalidEncodedHash
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, err
	}
	if version != argon2.Version {
		return nil, ErrIncompatibleArgon2Version
	}

	_, err = fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&result.Params.Memory,
		&result.Params.Iterations,
		&result.Params.Parallelism,
	)
	if err != nil {
		return nil, err
	}

	result.Salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, err
	}

	result.Hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, err
	}
	result.Params.KeyLength = uint32(len(result.Hash))

	return &result, nil
}
