package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/juancwu/secrething/internal/server/config"
	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/utils"
	"github.com/sumup/typeid"
)

type PackageType string

const (
	StdPackage PackageType = "std"
	UrlPackage PackageType = "url"
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
	// TokenTypeAccountActivate is for verifying a new user email and activate the account
	TokenTypeAccountActivate = "account_activate"
)

// Token duration constants
const (
	AccessTokenDuration          = 15 * time.Minute
	RefreshTokenDuration         = 7 * 24 * time.Hour // 7 days
	TempTokenDuration            = 5 * time.Minute
	APITokenDuration             = 90 * 24 * time.Hour // 90 days
	AccountActivateTokenDuration = 24 * time.Hour
)

const (
	ClientTypeWeb = "web"
	ClientTypeCLI = "cli"
	ClientTypeAPI = "api"
)

// Error type constants
const (
	TokenServiceErrExpired            = "token_service_expired"
	TokenServiceErrInvalid            = "token_service_invalid"
	TokenServiceErrGeneration         = "token_service_generation"
	TokenServiceErrEncryption         = "token_service_encryption"
	TokenServiceErrDecryption         = "token_service_decryption"
	TokenServiceErrDatabase           = "token_service_database"
	TokenServiceErrInvalidPackageType = "token_service_invalid_package_type"
)

type TokenServiceError struct {
	Err      error
	Template string
	Type     string
	Params   map[string]interface{}
}

func NewTokenServiceError(template, typ string, err error, params map[string]interface{}) TokenServiceError {
	if params == nil {
		params = make(map[string]interface{})
	}

	_, exists := params["err"]
	if !exists && err != nil {
		params["err"] = err
	}

	return TokenServiceError{
		Err:      err,
		Template: template,
		Type:     typ,
		Params:   params,
	}
}

// Error implements the error interface for TokenServiceError.
func (e TokenServiceError) Error() string {
	return utils.Interpolate(e.Template, e.Params)
}

// IsType is a utility function that checks if the error is of the given type.
func (e TokenServiceError) IsType(errType string) bool {
	return e.Type == errType
}

// Unwrap returns the underlying error
func (e TokenServiceError) Unwrap() error {
	return e.Err
}

// NewTokenExpiredError returns a new error for an expired token
func NewTokenExpiredError(expiryTime string) TokenServiceError {
	return NewTokenServiceError(
		"Token expired. Expected expiry time: {expiry_time}",
		TokenServiceErrExpired,
		nil,
		map[string]interface{}{"expiry_time": expiryTime},
	)
}

// NewTokenInvalidError returns a new error for an invalid token
func NewTokenInvalidError(err error, details string) TokenServiceError {
	return NewTokenServiceError(
		"Invalid token: {details}",
		TokenServiceErrInvalid,
		err,
		map[string]interface{}{"details": details},
	)
}

// NewTokenGenerationError returns a new error for token generation failures
func NewTokenGenerationError(err error) TokenServiceError {
	return NewTokenServiceError(
		"Failed to generate token: {err}",
		TokenServiceErrGeneration,
		err,
		nil,
	)
}

// NewTokenDatabaseError returns a new error for database operations related to tokens
func NewTokenDatabaseError(err error, operation string) TokenServiceError {
	return NewTokenServiceError(
		"Database error during {operation}: {err}",
		TokenServiceErrDatabase,
		err,
		map[string]interface{}{"operation": operation},
	)
}

// NewTokenInvalidPackageTypeError returns a new error for invalid package type given when generating tokens
func NewTokenInvalidPackageTypeError(err error, operation string) TokenServiceError {
	return NewTokenServiceError(
		"Invalid package type error during {operation}: {err}",
		TokenServiceErrInvalidPackageType,
		err,
		map[string]interface{}{"operation": operation},
	)
}

// TokenPayload represents the payload that is encrypted and is sent back to the user.
// This payload helps the server later verify the integrity and authenticity of the token
// since it is encrypted using AES256-GCM.
type TokenPayload struct {
	TokenID    db.TokenID `json:"tid"`
	TokenType  string     `json:"ttp"`
	ClientType string     `json:"ctp"`
	ExpiresAt  string     `json:"exp"`
	UserID     db.UserID  `json:"uid"`
}

// TokenPair represents a pair of access and refresh tokens.
// The tokens are encrypted TokenPayload in their string representation
// using base64 URL safe encoding.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

// GetTokenByID gets a token from the database that matches the given tokenID
func (s *TokenService) GetTokenByID(ctx context.Context, tokenID db.TokenID) (*db.Token, error) {
	q, err := db.Query()
	if err != nil {
		return nil, err
	}

	token, err := q.GetTokenByID(ctx, tokenID)
	return &token, err
}

// VerifyToken does a basic db lookup for the token and matches the token type and
// makes sure the token hasn't expired yet.
func (s *TokenService) VerifyToken(ctx context.Context, tokenStr string, packageType PackageType) (payload *TokenPayload, err error) {
	payload, err = s.Unpack(tokenStr, packageType)
	if err != nil {
		return nil, NewTokenInvalidError(err, "Failed to unpack token")
	}

	expiryTime, err := utils.ParseRFC3339NanoStr(payload.ExpiresAt)
	if err != nil {
		return nil, NewTokenInvalidError(err, "Invalid expiry time format")
	}

	if time.Now().After(expiryTime) {
		return nil, NewTokenExpiredError(payload.ExpiresAt)
	}

	return payload, nil
}

func (s *TokenService) GenerateTokenPair(ctx context.Context, userID db.UserID, clientType string) (TokenPair, error) {
	now := time.Now()

	// Perform actions in transaction to ensure both token types are made successfully.
	tx, q, err := db.QueryWithTx()
	if err != nil {
		return TokenPair{}, NewTokenDatabaseError(err, "begin transaction")
	}
	defer tx.Rollback()

	refreshToken, err := s.generateTokenPayload(ctx, q, userID, TokenTypeRefresh, clientType, now, now.Add(RefreshTokenDuration))
	if err != nil {
		return TokenPair{}, err // generateToken already returns TokenServiceError
	}

	accessToken, err := s.generateTokenPayload(ctx, q, userID, TokenTypeAccess, clientType, now, now.Add(AccessTokenDuration))
	if err != nil {
		return TokenPair{}, err // generateToken already returns TokenServiceError
	}

	rtk, err := s.StdPack(refreshToken)
	if err != nil {
		return TokenPair{}, NewTokenServiceError(
			"Failed to encrypt refresh token: {err}",
			TokenServiceErrEncryption,
			err,
			nil,
		)
	}

	atk, err := s.StdPack(accessToken)
	if err != nil {
		return TokenPair{}, NewTokenServiceError(
			"Failed to encrypt access token: {err}",
			TokenServiceErrEncryption,
			err,
			nil,
		)
	}

	// Commit changes only when everything has gone well and the tokens can be sent back to the client
	if err := tx.Commit(); err != nil {
		return TokenPair{}, NewTokenDatabaseError(err, "commit transaction")
	}

	return TokenPair{
		RefreshToken: rtk,
		AccessToken:  atk,
	}, nil
}

func (s *TokenService) NewAccountActivateToken(ctx context.Context, userID db.UserID) (string, error) {
	now := time.Now()
	exp := now.Add(AccountActivateTokenDuration)
	return s.generateToken(
		ctx,
		UrlPackage,
		userID,
		TokenTypeAccountActivate,
		"?",
		now,
		exp,
	)
}

func (s *TokenService) DecryptToken(data []byte) ([]byte, error) {
	decrypted, err := utils.DecryptAES(data, s.key)
	if err != nil {
		return nil, NewTokenServiceError(
			"Failed to decrypt token: {err}",
			TokenServiceErrDecryption,
			err,
			nil,
		)
	}
	return decrypted, nil
}

func (s *TokenService) EncryptToken(data []byte) ([]byte, error) {
	encrypted, err := utils.EncryptAES(data, s.key)
	if err != nil {
		return nil, NewTokenServiceError(
			"Failed to encrypt token: {err}",
			TokenServiceErrEncryption,
			err,
			nil,
		)
	}
	return encrypted, nil
}

func (s *TokenService) Unpack(t string, packType PackageType) (*TokenPayload, error) {
	switch packType {
	case StdPackage:
		return s.StdUnpack(t)
	case UrlPackage:
		return s.UrlUnpack(t)
	}
	return nil, NewTokenInvalidPackageTypeError(fmt.Errorf("unknown package type: %s", packType), "token unpacking")
}

func (s *TokenService) StdPack(t *TokenPayload) (string, error) {
	data, err := s.pack(t)
	if err != nil {
		return "", err // pack already returns TokenServiceError
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *TokenService) UrlPack(t *TokenPayload) (string, error) {
	data, err := s.pack(t)
	if err != nil {
		return "", err // pack already returns TokenServiceError
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func (s *TokenService) StdUnpack(t string) (*TokenPayload, error) {
	b, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		return nil, NewTokenInvalidError(err, "Invalid base64 encoding")
	}
	return s.unpack(b)
}

func (s *TokenService) UrlUnpack(t string) (*TokenPayload, error) {
	b, err := base64.URLEncoding.DecodeString(t)
	if err != nil {
		return nil, NewTokenInvalidError(err, "Invalid URL-safe base64 encoding")
	}
	return s.unpack(b)
}

func (s *TokenService) generateToken(ctx context.Context, packageType PackageType, userID db.UserID, tokenType, clientType string, now, exp time.Time) (string, error) {
	var packaged string
	var err error
	var tx *sql.Tx
	var q *db.Queries

	// Perform actions in transaction to ensure both token types are made successfully.
	tx, q, err = db.QueryWithTx()
	if err != nil {
		return "", NewTokenDatabaseError(err, "begin transaction")
	}
	defer tx.Rollback()

	token, err := s.generateTokenPayload(
		ctx,
		q,
		userID,
		tokenType,
		clientType,
		now,
		exp,
	)
	if err != nil {
		return "", err // generateToken already returns TokenServiceError
	}

	switch packageType {
	case StdPackage:
		packaged, err = s.StdPack(token)
	case UrlPackage:
		packaged, err = s.UrlPack(token)
	default:
		return "", NewTokenInvalidPackageTypeError(fmt.Errorf("unknown package type: %s", packageType), "token packaging")
	}
	if err != nil {
		return "", err // Pack methods already return a TokenServiceError
	}

	if err := tx.Commit(); err != nil {
		return "", NewTokenDatabaseError(err, "commit transaction")
	}

	return packaged, nil
}

func (s *TokenService) generateTokenPayload(ctx context.Context, q *db.Queries, userID db.UserID, tokenType, clientType string, now time.Time, exp time.Time) (*TokenPayload, error) {
	tokenID, err := s.generateTokenID()
	if err != nil {
		// GenerateTokenID already returns TokenServiceError
		return nil, err
	}

	token, err := q.CreateToken(ctx, db.CreateTokenParams{
		TokenID:    tokenID,
		UserID:     userID,
		TokenType:  tokenType,
		ClientType: clientType,
		CreatedAt:  utils.FormatRFC3339NanoFixed(now),
		ExpiresAt:  utils.FormatRFC3339NanoFixed(exp),
	})
	if err != nil {
		return nil, NewTokenDatabaseError(err, "create token")
	}

	return &TokenPayload{
		TokenID:    token.TokenID,
		TokenType:  token.TokenType,
		ClientType: token.ClientType,
		UserID:     token.UserID,
		ExpiresAt:  token.ExpiresAt,
	}, nil
}

// generateTokenID creates a unique token identifier
func (*TokenService) generateTokenID() (db.TokenID, error) {
	// Generate a random token ID
	tokenID, err := typeid.New[db.TokenID]()
	if err != nil {
		return db.TokenID{}, NewTokenGenerationError(err)
	}
	return tokenID, nil
}

func (s *TokenService) pack(t *TokenPayload) ([]byte, error) {
	raw, err := json.Marshal(t)
	if err != nil {
		return nil, NewTokenServiceError(
			"Failed to marshal token payload: {err}",
			TokenServiceErrEncryption,
			err,
			nil,
		)
	}
	enc, err := s.EncryptToken(raw)
	if err != nil {
		// EncryptToken already returns TokenServiceError
		return nil, err
	}
	return enc, nil
}

func (s *TokenService) unpack(data []byte) (*TokenPayload, error) {
	d, err := s.DecryptToken(data)
	if err != nil {
		// DecryptToken already returns TokenServiceError
		return nil, err
	}
	var payload TokenPayload
	if err := json.Unmarshal(d, &payload); err != nil {
		return nil, NewTokenInvalidError(err, "Invalid token format")
	}
	return &payload, nil
}
