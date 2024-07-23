package jwt

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/juancwu/konbini/store"
)

const (
	// matches the access enum (TOKEN) from the db
	ACCESS_TOKEN_DB_ENUM = "access"
	// matches the refresh enum (TOKEN) from the db
	REFRESH_TOKEN_DB_ENUM = "refresh"
)

// JwtClaims represents the body of jwt tokens generated when a user signins.
type JwtClaims struct {
	// UserId for easy identification of the user
	UserId string `json:"user_id"`
	// TokenType for a way to differentiate access and refresh tokens
	TokenType string `json:"token_type"`
	// jwt.RegisteredClaims are the default normal claims that a jwt should have
	jwt.RegisteredClaims
}

// Creates a new row in the auth_tokens table to register the token
// to allow validation/invalidation of tokens.
func saveTokenInDatabase(tx *sql.Tx, uid, tokenType string, exp time.Time) (string, error) {
	row := tx.QueryRow("INSERT INTO auth_tokens (user_id, token_type, expires_at) VALUES ($1, $2, $3) RETURNING id;", uid, tokenType, exp)
	if err := row.Err(); err != nil {
		return "", err
	}
	var oid string
	if err := row.Scan(&oid); err != nil {
		return "", err
	}
	return oid, nil
}

// GenerateAccessToken generates an access token.
// Access tokens have 10 minutes of lifespan.
// This uses a transaction to insert the generated token in
// the database so that validation and invalidation can take place.
//
// You must call tx.Commit for the changes to take effect and the token to be valid.
func GenerateAccessToken(tx *sql.Tx, uid string) (string, error) {
	exp := time.Now().Add(time.Minute * 10)
	oid, err := saveTokenInDatabase(tx, uid, os.Getenv("JWT_ACCESS_TOKEN_TYPE"), exp)
	if err != nil {
		return "", err
	}
	return generateToken(uid, os.Getenv("JWT_ACCESS_TOKEN_TYPE"), os.Getenv("JWT_ACCESS_TOKEN_SECRET"), oid, exp)
}

// GenerateRefreshToken generates a refresh token.
// Refresh tokens have 1 week of lifepspan.
// This uses a transaction to insert the generated token in
// the database so that validation and invalidation can take place.
//
// You must call tx.Commit for the changes to take effect and the token to be valid.
func GenerateRefreshToken(tx *sql.Tx, uid string) (string, error) {
	exp := time.Now().Add(time.Hour * 24 * 7)
	oid, err := saveTokenInDatabase(tx, uid, os.Getenv("JWT_REFRESH_TOKEN_TYPE"), exp)
	if err != nil {
		return "", err
	}
	return generateToken(uid, os.Getenv("JWT_REFRESH_TOKEN_TYPE"), os.Getenv("JWT_REFRESH_TOKEN_SECRET"), oid, exp)
}

// VerifyAccessToken verifies the given access token.
func VerifyAccessToken(token string) (*JwtClaims, error) {
	claims := JwtClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if err := store.IsTokenValid(claims.ID, claims.TokenType); err != nil {
		return nil, err
	}
	return &claims, nil
}

// VerifyRefreshToken verifies the given refresh token.
func VerifyRefreshToken(token string) (*JwtClaims, error) {
	claims := JwtClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	// check with db
	if err := store.IsTokenValid(claims.ID, claims.TokenType); err != nil {
		return nil, err
	}
	return &claims, nil
}

// generateToken is a private method that generates and signs a jwt.
// Each type of token have their own secret that should be passed.
func generateToken(uid string, tokenType string, secret string, objectId string, exp time.Time) (string, error) {
	claims := JwtClaims{
		UserId:    uid,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("JWT_ISSUER"),
			ID:        objectId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
