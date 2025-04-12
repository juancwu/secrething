package services

import (
	"context"
	"sync"
	"time"

	"github.com/juancwu/secrething/internal/server/config"
	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/utils"
	"github.com/sumup/typeid"
)

// TokenType constants
const (
	// TokenTypeAccess is for short-lived access (for all clients)
	TokenTypeAccess = "access"
	// TokenTypeRefresh is for long-lived refresh (for all clients)
	TokenTypeRefresh = "refresh"
	// TokenTypeTemp is for TOTP verification (temporary)
	TokenTypeTemp = "temp"
	// TokenTypeAPI is for API authentication
	TokenTypeAPI = "api"
	// TokenAccountActivate is for verifying a new user email and activate the account
	TokenAccountActivate = "account_activate"
)

// Token duration constants
const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 days
	TempTokenDuration    = 5 * time.Minute
	APITokenDuration     = 90 * 24 * time.Hour // 90 days
)

const (
	ClientTypeWeb = "web"
	ClientTypeCLI = "cli"
	ClientTypeAPI = "api"
)

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  db.TokenID `json:"access_token"`
	RefreshToken db.TokenID `json:"refresh_token"`
}

// TokenService handles token generation, validation, and revocation
type TokenService struct {
	key []byte
}

var tokenService *TokenService
var tokenServiceMut sync.RWMutex

// NewTokenService creates a new token service with the specified key
func NewTokenService() *TokenService {
	tokenServiceMut.Lock()
	defer tokenServiceMut.Unlock()

	if tokenService == nil {
		// Convert the key from string to bytes
		key := []byte(config.Token().AuthKey)
		tokenService = &TokenService{key: key}
	}
	return tokenService
}

// GenerateTokenID creates a unique token identifier
func (*TokenService) GenerateTokenID() (db.TokenID, error) {
	// Generate a random token ID
	return typeid.New[db.TokenID]()
}

func (s *TokenService) GenerateTokenPair(ctx context.Context, userID db.UserID, clientType string) (TokenPair, error) {
	now := time.Now()

	refreshToken, err := s.GenerateToken(ctx, userID, TokenTypeRefresh, clientType, now, now.Add(RefreshTokenDuration))
	if err != nil {
		return TokenPair{}, err
	}

	accessToken, err := s.GenerateToken(ctx, userID, TokenTypeAccess, clientType, now, now.Add(AccessTokenDuration))
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		RefreshToken: refreshToken.TokenID,
		AccessToken:  accessToken.TokenID,
	}, nil
}

func (s *TokenService) GenerateToken(ctx context.Context, userID db.UserID, tokenType, clientType string, now time.Time, exp time.Time) (db.Token, error) {
	q, err := db.Query()
	if err != nil {
		return db.Token{}, err
	}

	tokenID, err := s.GenerateTokenID()
	if err != nil {
		return db.Token{}, err
	}

	return q.CreateToken(ctx, db.CreateTokenParams{
		TokenID:    tokenID,
		UserID:     userID,
		TokenType:  tokenType,
		ClientType: clientType,
		CreatedAt:  utils.FormatRFC3339NanoFixed(now),
		ExpiresAt:  utils.FormatRFC3339NanoFixed(exp),
	})
}
