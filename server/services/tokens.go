package services

import (
	"encoding/base64"
	"errors"
	"konbini/server/config"
	"konbini/server/utils"
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	FULL_USER_TOKEN_TYPE    TokenType = "full_token"
	PARTIAL_USER_TOKEN_TYPE TokenType = "partial_token"
	EMAIL_TOKEN_TYPE        string    = "email_token"

	customer string = "customer"
)

func (t TokenType) Valid() bool {
	return t == PARTIAL_USER_TOKEN_TYPE || t == FULL_USER_TOKEN_TYPE
}

func (t TokenType) String() string {
	if t == FULL_USER_TOKEN_TYPE {
		return "full_token"
	}
	return "partial_token"
}

var (
	ErrInvalidTokenType    error = errors.New("Invalid token type. Use constants to not make a mistake.")
	ErrExpiredJWT          error = errors.New("AuthToken has expired.")
	ErrInvalidAuthTokenLen error = errors.New("Invalid token length (>102)")
)

type AuthToken struct {
	ID        string
	UserID    string
	TokenType TokenType
	ExpiresAt time.Time
}

func (t *AuthToken) Package() (string, error) {
	cfg, err := config.Global()
	if err != nil {
		return "", err
	}

	id := []byte(t.ID)
	userID := []byte(t.UserID)
	expiresAt := []byte(utils.FormatRFC3339NanoFixed(t.ExpiresAt))
	tokenType := []byte(t.TokenType)

	// id + userID + expiresAt + tokenType
	// 36 + 36 + 30 + len(tokenType)
	data := make([]byte, 102+len(tokenType))
	copy(data[0:], id)
	copy(data[36:], userID)
	copy(data[72:], expiresAt)
	copy(data[102:], []byte(tokenType))

	ciphertext, err := utils.EncryptAES(data, cfg.GetFullTokenKey())
	if err != nil {
		return "", err
	}

	// encode in base64
	b64Cipher := base64.URLEncoding.EncodeToString(ciphertext)

	return b64Cipher, nil
}

func NewAuthToken(id, userID string, tokenType TokenType, exp time.Time) (*AuthToken, error) {
	if !tokenType.Valid() {
		return nil, ErrInvalidTokenType
	}

	j := &AuthToken{
		ID:        id,
		UserID:    userID,
		TokenType: tokenType,
		ExpiresAt: exp,
	}

	return j, nil
}

func VerifyAuthToken(token string) (*AuthToken, error) {
	cfg, err := config.Global()
	if err != nil {
		return nil, err
	}

	plaintext, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	plaintext, err = utils.DecryptAES(plaintext, cfg.GetFullTokenKey())
	if err != nil {
		return nil, err
	}

	if len(plaintext) <= 102 {
		return nil, ErrInvalidAuthTokenLen
	}

	id := plaintext[0:36]
	userID := plaintext[36:72]
	expiresAtBytes := plaintext[72:102]
	expiresAt, err := time.Parse(time.RFC3339Nano, string(expiresAtBytes))
	if err != nil {
		return nil, err
	}
	tokenType := TokenType(string(plaintext[102:]))
	if !tokenType.Valid() {
		return nil, ErrInvalidTokenType
	}

	if time.Now().After(expiresAt) {
		return nil, ErrExpiredJWT
	}

	return &AuthToken{
		ID:        string(id),
		UserID:    string(userID),
		TokenType: tokenType,
		ExpiresAt: expiresAt.UTC(),
	}, nil
}

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
