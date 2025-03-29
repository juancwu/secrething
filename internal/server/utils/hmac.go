package utils

import (
	"crypto/hmac"
	"crypto/sha256"
)

func CreateHMAC(message, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	hash := h.Sum(nil)
	return hash
}

func VerifyHMAC(message, key, expectedMAC []byte) bool {
	mac := CreateHMAC(message, key)
	return hmac.Equal(mac, expectedMAC)
}
