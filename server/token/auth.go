package token

import (
	"encoding/base64"
	"encoding/json"

	"github.com/juancwu/konbini/server/utils"
)

type AuthTokenType string

const (
	TemporaryToken      AuthTokenType = "temporary_token"
	LimitedAccessToken  AuthTokenType = "limited_access_token"
	AccessToken         AuthTokenType = "access_token"
	LimitedRefreshToken AuthTokenType = "limited_refresh_token"
	RefreshToken        AuthTokenType = "refresh_token"
)

func (t AuthTokenType) IsValid() bool {
	switch t {
	case TemporaryToken, AccessToken, RefreshToken:
		return true
	}
	return false
}

type AuthToken struct {
	UserID string
	Type   AuthTokenType
}

// NewAuthToken creates a new auth token with the given token type.
func NewAuthToken(userID string, tokType AuthTokenType) AuthToken {
	return AuthToken{
		UserID: userID,
		Type:   tokType,
	}
}

// Scan reads the data given and store it in the current AuthToken.
func (t *AuthToken) Scan(data []byte) error {
	return json.Unmarshal(data, t)
}

// Package encrypts and encodes (base64-url) the AuthToken which makes it ready for network transfer.
func (t AuthToken) Package(key []byte) (string, error) {
	// Marshal token
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	// Encrypt token
	encrypted, err := utils.EncryptAES(data, key)
	if err != nil {
		return "", err
	}

	// Encode in base64 URL safe
	b64Token := base64.URLEncoding.EncodeToString(encrypted)

	return b64Token, nil
}

// NewTempAuthToken creates a new temporary token.
func NewTemporaryToken(userID string) AuthToken {
	return NewAuthToken(userID, TemporaryToken)
}

// NewAccessToken creates a new access token.
func NewAccessToken(userID string) AuthToken {
	return NewAuthToken(userID, AccessToken)
}

// NewRefreshToken creates a new refresh token.
func NewRefreshToken(userID string) AuthToken {
	return NewAuthToken(userID, RefreshToken)
}

// NewLimitedAccessToken creates a new limited access token.
func NewLimitedAccessToken(userID string) AuthToken {
	return NewAuthToken(userID, LimitedAccessToken)
}

// NewLimitedRefreshToken creates a new limited refresh token
func NewLimitedRefreshToken(userID string) AuthToken {
	return NewAuthToken(userID, LimitedRefreshToken)
}
