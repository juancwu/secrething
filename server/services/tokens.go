package services

import (
	"encoding/base64"
	"errors"
	"konbini/server/config"
	"konbini/server/utils"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TokenType byte

const (
	FULL_USER_TOKEN_TYPE    TokenType = 1
	PARTIAL_USER_TOKEN_TYPE TokenType = 0
	EMAIL_TOKEN_TYPE        string    = "email_token"

	customer string = "customer"
)

func TokenTypeFromByte(b byte) (TokenType, error) {
	t := TokenType(b)
	if !t.Valid() {
		return 0, ErrInvalidTokenType
	}
	return t, nil
}

func TokenTypeFromString(s string) (TokenType, error) {
	switch s {
	case "full_token":
		return FULL_USER_TOKEN_TYPE, nil
	case "partial_token":
		return PARTIAL_USER_TOKEN_TYPE, nil
	}
	return TokenType(0), ErrInvalidTokenType
}

func (t TokenType) Valid() bool {
	return t == PARTIAL_USER_TOKEN_TYPE || t == FULL_USER_TOKEN_TYPE
}

func (t TokenType) Byte() byte {
	if t == FULL_USER_TOKEN_TYPE {
		return 1
	}
	return 0
}

func (t TokenType) String() string {
	if t == FULL_USER_TOKEN_TYPE {
		return "full_token"
	}
	return "partial_token"
}

var (
	ErrInvalidTokenType error = errors.New("Invalid token type. Use constants to not make a mistake.")
	ErrExpiredJWT       error = errors.New("AuthToken has expired.")
)

type AuthToken struct {
	ID        string
	UserID    string
	TokenType TokenType
	ExpiresAt time.Time
}

func (j *AuthToken) EncryptedString() (string, error) {
	cfg, err := config.Global()
	if err != nil {
		return "", err
	}

	id := []byte(j.ID)
	userID := []byte(j.UserID)
	expiresAt := []byte(j.ExpiresAt.Format(time.RFC3339))

	// 1 for the token type
	data := make([]byte, len(id)+len(userID)+len(expiresAt)+1)
	data[0] = j.TokenType.Byte()
	offset := 1
	copy(data[offset:], id)
	offset += len(id)
	copy(data[offset:], userID)
	offset += len(userID)
	copy(data[offset:], expiresAt)

	ciphertext, err := utils.EncryptAES(data, cfg.GetFullTokenKey())
	if err != nil {
		return "", err
	}

	log.Info().Bytes("ciphertext", ciphertext).Msg("debug")

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

	tokenType, err := TokenTypeFromByte(plaintext[0])
	if err != nil {
		return nil, err
	}

	offset := 1
	uuidv4Len := 36
	id := plaintext[offset : offset+uuidv4Len]
	offset += uuidv4Len
	userID := plaintext[offset : offset+uuidv4Len]
	offset += uuidv4Len
	expiresAtBytes := plaintext[offset:]
	expiresAt, err := time.Parse(time.RFC3339, string(expiresAtBytes))
	if err != nil {
		return nil, err
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
