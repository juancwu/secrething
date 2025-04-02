package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/juancwu/secrething/internal/server/config"
	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/utils"
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
)

// Token duration constants
const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 days
	TempTokenDuration    = 5 * time.Minute
	APITokenDuration     = 90 * 24 * time.Hour // 90 days
)

// ClientType identifies the type of client that is using the token
type ClientType string

const (
	ClientTypeWeb ClientType = "web"
	ClientTypeCLI ClientType = "cli"
	ClientTypeAPI ClientType = "api"
)

// TokenPayload contains the minimal user information
// needed for authentication
type TokenPayload struct {
	// UserID is the unique identifier for the user
	UserID string `json:"uid"`
	// TokenID is the unique identifier for this token
	TokenID string `json:"tid"`
	// TokenType indicates what type of token this is
	TokenType string `json:"typ"`
	// ClientType indicates what client type is using this token
	ClientType ClientType `json:"cli"`
	// ExpiresAt is the expiration time for this token
	ExpiresAt time.Time `json:"exp"`
	// RequiresTotp indicates if this token requires TOTP verification
	RequiresTotp bool `json:"totp,omitempty"`
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"` // omitted for web clients
	ExpiresIn    int    `json:"expires_in"`              // access token lifetime in seconds
}

// TempTokenResponse represents a temporary token for TOTP verification
type TempTokenResponse struct {
	TempToken string `json:"temp_token"`
	ExpiresIn int    `json:"expires_in"` // temporary token lifetime in seconds
}

// TokenService handles token generation, validation, and revocation
type TokenService struct {
	key []byte
}

// NewTokenService creates a new token service with the specified key
func NewTokenService() *TokenService {
	// Convert the key from string to bytes
	key := []byte(config.Token().AuthKey)
	return &TokenService{key: key}
}

// GenerateTokenID creates a unique token identifier
func GenerateTokenID() (string, error) {
	// Generate a random token ID
	tokenIDBytes := make([]byte, 16)
	if _, err := rand.Read(tokenIDBytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tokenIDBytes), nil
}

// CreateToken creates a token with the given payload
func (s *TokenService) CreateToken(payload TokenPayload) (string, error) {
	// Serialize payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Encrypt the payload
	encrypted, err := utils.EncryptAES(jsonPayload, s.key)
	if err != nil {
		return "", err
	}

	// Base64 encode the encrypted payload to make it URL-safe
	encodedToken := base64.RawURLEncoding.EncodeToString(encrypted)

	// Add a prefix to identify the token type
	return fmt.Sprintf("%s_%s", payload.TokenType, encodedToken), nil
}

// GenerateTokenPair creates an access and refresh token pair for the specified user
func (s *TokenService) GenerateTokenPair(ctx context.Context, userID string, clientType ClientType) (*TokenPair, error) {
	// Step 1: Create the access token
	accessTokenID, err := GenerateTokenID()
	if err != nil {
		return nil, err
	}

	accessPayload := TokenPayload{
		UserID:     userID,
		TokenID:    accessTokenID,
		TokenType:  TokenTypeAccess,
		ClientType: clientType,
		ExpiresAt:  time.Now().Add(AccessTokenDuration),
	}

	accessToken, err := s.CreateToken(accessPayload)
	if err != nil {
		return nil, err
	}

	// Step 2: Create the refresh token
	refreshTokenID, err := GenerateTokenID()
	if err != nil {
		return nil, err
	}

	refreshPayload := TokenPayload{
		UserID:     userID,
		TokenID:    refreshTokenID,
		TokenType:  TokenTypeRefresh,
		ClientType: clientType,
		ExpiresAt:  time.Now().Add(RefreshTokenDuration),
	}

	refreshToken, err := s.CreateToken(refreshPayload)
	if err != nil {
		return nil, err
	}

	// Store the refresh token in the database
	q, err := db.Query()
	if err != nil {
		return nil, err
	}

	// Store refresh token with client type to differentiate between platforms
	dbTokenType := string(TokenTypeRefresh + "_" + string(clientType))
	_, err = q.CreateUserToken(ctx, db.CreateUserTokenParams{
		UserID:    userID,
		TokenType: dbTokenType,
		ExpiresAt: utils.FormatRFC3339NanoFixed(refreshPayload.ExpiresAt),
		CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	// Create the response based on client type
	pair := &TokenPair{
		AccessToken: accessToken,
		ExpiresIn:   int(AccessTokenDuration.Seconds()),
	}

	// For non-web clients, include the refresh token in the response
	// Web clients should get the refresh token via HTTP-only cookie
	if clientType != ClientTypeWeb {
		pair.RefreshToken = refreshToken
	}

	return pair, nil
}

// GenerateAPIToken creates a long-lived API token for admin access
func (s *TokenService) GenerateAPIToken(ctx context.Context, userID string) (string, error) {
	tokenID, err := GenerateTokenID()
	if err != nil {
		return "", err
	}

	payload := TokenPayload{
		UserID:     userID,
		TokenID:    tokenID,
		TokenType:  TokenTypeAPI,
		ClientType: ClientTypeAPI,
		ExpiresAt:  time.Now().Add(APITokenDuration),
	}

	token, err := s.CreateToken(payload)
	if err != nil {
		return "", err
	}

	// Store the API token in the database
	q, err := db.Query()
	if err != nil {
		return "", err
	}

	_, err = q.CreateUserToken(ctx, db.CreateUserTokenParams{
		UserID:    userID,
		TokenType: TokenTypeAPI,
		ExpiresAt: utils.FormatRFC3339NanoFixed(payload.ExpiresAt),
		CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

// GenerateTempToken creates a temporary token for TOTP verification
func (s *TokenService) GenerateTempToken(ctx context.Context, userID string) (*TempTokenResponse, error) {
	tokenID, err := GenerateTokenID()
	if err != nil {
		return nil, err
	}

	// Create token payload with TOTP flag
	payload := TokenPayload{
		UserID:       userID,
		TokenID:      tokenID,
		TokenType:    TokenTypeTemp,
		ClientType:   ClientTypeWeb, // TOTP verification is typically done via web
		ExpiresAt:    time.Now().Add(TempTokenDuration),
		RequiresTotp: true,
	}

	token, err := s.CreateToken(payload)
	if err != nil {
		return nil, err
	}

	// Store the token in the database
	q, err := db.Query()
	if err != nil {
		return nil, err
	}

	_, err = q.CreateUserToken(ctx, db.CreateUserTokenParams{
		UserID:    userID,
		TokenType: TokenTypeTemp,
		ExpiresAt: utils.FormatRFC3339NanoFixed(payload.ExpiresAt),
		CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	return &TempTokenResponse{
		TempToken: token,
		ExpiresIn: int(TempTokenDuration.Seconds()),
	}, nil
}

// RefreshTokens validates a refresh token and returns a new token pair
func (s *TokenService) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate the refresh token
	payload, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Ensure it's a refresh token
	if payload.TokenType != TokenTypeRefresh {
		return nil, errors.New("not a refresh token")
	}

	// Check if the token is expired
	if time.Now().After(payload.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	// Verify the token in the database
	q, err := db.Query()
	if err != nil {
		return nil, err
	}

	// Check if this token type exists for this user
	dbTokenType := string(TokenTypeRefresh + "_" + string(payload.ClientType))
	_, err = q.GetUserTokenByType(ctx, db.GetUserTokenByTypeParams{
		UserID:    payload.UserID,
		TokenType: dbTokenType,
	})
	if err != nil {
		if db.IsNoRows(err) {
			return nil, errors.New("token not found or revoked")
		}
		return nil, err
	}

	// Generate a new token pair
	return s.GenerateTokenPair(ctx, payload.UserID, payload.ClientType)
}

// ParseToken extracts and verifies the token payload without database check
func (s *TokenService) ParseToken(token string) (*TokenPayload, error) {
	// Split the token into type and payload
	parts := strings.SplitN(token, "_", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	tokenType := parts[0]
	encodedPayload := parts[1]

	// Decode the base64 encoded payload
	encrypted, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return nil, err
	}

	// Decrypt the payload
	decrypted, err := utils.DecryptAES(encrypted, s.key)
	if err != nil {
		return nil, err
	}

	// Deserialize the JSON payload
	var payload TokenPayload
	if err := json.Unmarshal(decrypted, &payload); err != nil {
		return nil, err
	}

	// Verify token type matches the prefix
	if payload.TokenType != tokenType {
		return nil, errors.New("token type mismatch")
	}

	return &payload, nil
}

// ValidateToken validates a token and returns the payload if valid
func (s *TokenService) ValidateToken(ctx context.Context, token string) (*TokenPayload, error) {
	// Parse the token first
	payload, err := s.ParseToken(token)
	if err != nil {
		return nil, err
	}

	// Check if the token is expired
	if time.Now().After(payload.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	// For access tokens, we don't need to check the database
	// They're short-lived and validated cryptographically
	if payload.TokenType == TokenTypeAccess {
		return payload, nil
	}

	// For other token types, check the database
	q, err := db.Query()
	if err != nil {
		return nil, err
	}

	// Determine the token type for database query
	dbTokenType := payload.TokenType
	if payload.TokenType == TokenTypeRefresh {
		// For refresh tokens, we store them with client type
		dbTokenType = string(TokenTypeRefresh + "_" + string(payload.ClientType))
	}

	// Check if this token type exists for this user
	_, err = q.GetUserTokenByType(ctx, db.GetUserTokenByTypeParams{
		UserID:    payload.UserID,
		TokenType: dbTokenType,
	})
	if err != nil {
		if db.IsNoRows(err) {
			return nil, errors.New("token not found or revoked")
		}
		return nil, err
	}

	return payload, nil
}

// RevokeToken revokes a specific refresh or API token
func (s *TokenService) RevokeToken(ctx context.Context, token string) error {
	// Parse the token
	payload, err := s.ParseToken(token)
	if err != nil {
		return err
	}

	// Only refresh tokens and API tokens can be revoked
	if payload.TokenType != TokenTypeRefresh && payload.TokenType != TokenTypeAPI {
		return errors.New("only refresh and API tokens can be revoked")
	}

	q, err := db.Query()
	if err != nil {
		return err
	}

	// Determine the token type for database
	dbTokenType := payload.TokenType
	if payload.TokenType == TokenTypeRefresh {
		dbTokenType = string(TokenTypeRefresh + "_" + string(payload.ClientType))
	}

	return q.DeleteUserToken(ctx, db.DeleteUserTokenParams{
		UserID:    payload.UserID,
		TokenType: dbTokenType,
	})
}

// RevokeAllTokens revokes all tokens for a user
func (s *TokenService) RevokeAllTokens(ctx context.Context, userID string) error {
	q, err := db.Query()
	if err != nil {
		return err
	}

	return q.DeleteAllUserTokens(ctx, userID)
}

// RevokeAllClientTokens revokes all tokens for a specific client type
func (s *TokenService) RevokeAllClientTokens(ctx context.Context, userID string, clientType ClientType) error {
	q, err := db.Query()
	if err != nil {
		return err
	}

	// For refresh tokens with client type
	tokenType := string(TokenTypeRefresh + "_" + string(clientType))

	return q.DeleteUserTokensByType(ctx, db.DeleteUserTokensByTypeParams{
		UserID:    userID,
		TokenType: tokenType,
	})
}
