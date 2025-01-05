package services

import (
	"encoding/base64"
	"konbini/server/config"
	"konbini/server/utils"
	"time"

	"github.com/google/uuid"
)

const (
	FULL_USER_TOKEN_TYPE    string = "full_user_token"
	PARTIAL_USER_TOKEN_TYPE string = "partial_user_token"
	EMAIL_TOKEN_TYPE        string = "email_token"
)

type EmailToken struct {
	Id        string
	UserId    string
	Hmac      []byte
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewEmailToken(userId string) (*EmailToken, error) {
	if _, err := uuid.Parse(userId); err != nil {
		return nil, err
	}
	emailTokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	exp := now.Add(10 * time.Minute)
	return &EmailToken{
		Id:        emailTokenId.String(),
		UserId:    userId,
		Hmac:      nil,
		CreatedAt: now,
		ExpiresAt: exp,
	}, nil
}

// ExtractEmailTokenId extracts the email token id from the encoded token and also verifies integrity.
// It is expected that the token is base64 encoded using the url standard.
func ExtractEmailTokenId(b64Token string) (string, error) {
	cfg, err := config.Global()
	if err != nil {
		return "", err
	}

	decodedToken, err := base64.URLEncoding.DecodeString(b64Token)
	if err != nil {
		return "", err
	}

	id, err := utils.DecryptAES(decodedToken, cfg.GetEmailTokenKey())
	if err != nil {
		return "", err
	}

	return string(id), nil
}

// Package packages the email token into a string that can be sent as a token.
func (t *EmailToken) Package() (string, error) {
	cfg, err := config.Global()
	if err != nil {
		return "", err
	}

	idBytes := []byte(t.Id)
	// encryption uses aes-gcm so it encrypts and authenticate message integrity
	encryptedId, err := utils.EncryptAES(idBytes, cfg.GetEmailTokenKey())
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(encryptedId), nil
}
