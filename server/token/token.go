package token

import (
	"encoding/base64"

	"github.com/juancwu/konbini/server/utils"
)

type Token interface {
	Package(key []byte) (string, error)
}

// Unpack is a universal method to decode and decrypt tokens.
func Unpack(b64Tok string, key []byte) ([]byte, error) {
	// Decode the base64-url encoded token
	decoded, err := base64.URLEncoding.DecodeString(b64Tok)
	if err != nil {
		return nil, err
	}

	// Decrypt the decoded token
	decrypted, err := utils.DecryptAES(decoded, key)
	return decrypted, err
}
