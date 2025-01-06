package services

import (
	"encoding/base64"
	"errors"
	"konbini/server/config"
	"konbini/server/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	FULL_USER_TOKEN_TYPE    string = "full_user_token"
	PARTIAL_USER_TOKEN_TYPE string = "partial_user_token"
	EMAIL_TOKEN_TYPE        string = "email_token"

	customer string = "customer"
)

var (
	ErrInvalidTokenType error = errors.New("Invalid token type. Use constants to not make a mistake.")
)

type JWTClaims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type JWT struct {
	Claims JWTClaims
}

func (j *JWT) SignedString() (string, error) {
	cfg, err := config.Global()
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j.Claims)
	key, err := getJWTKeyByType(j.Claims.Type, cfg)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}

// NewJWT generates a new JWT and signs it with HS256.
// An id must be provided since it is what relates the JWT to the
// row stored in the database.
func NewJWT(id, tokType string, expiresAt time.Time) (*JWT, error) {
	if !isValidJWTType(tokType) {
		return nil, ErrInvalidTokenType
	}

	claims := JWTClaims{
		Type: tokType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Konbini",
			ID:        id,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Audience:  []string{customer},
		},
	}
	j := &JWT{Claims: claims}
	return j, nil
}

// isValidJWTType checks if the given string is a valid JWT type string
func isValidJWTType(tokType string) bool {
	return tokType == PARTIAL_USER_TOKEN_TYPE || tokType == FULL_USER_TOKEN_TYPE
}

// getJWTKeyByType gets the correct key based on the given JWT type string.
// Returns an error if the token type is not valid.
func getJWTKeyByType(tokType string, cfg *config.Config) ([]byte, error) {
	switch tokType {
	case FULL_USER_TOKEN_TYPE:
		return cfg.GetFullTokenKey(), nil
	case PARTIAL_USER_TOKEN_TYPE:
		return cfg.GetPartialTokenKey(), nil
	}

	return nil, ErrInvalidTokenType
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
