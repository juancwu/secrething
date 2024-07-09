package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// GenerateAccessToken generates an access token.
// Access tokens have 10 minutes of lifespan.
func GenerateAccessToken(uid string) (string, error) {
	exp := time.Now().Add(time.Minute * 10)
	return generateToken(uid, os.Getenv("JWT_ACCESS_TOKEN_TYPE"), os.Getenv("JWT_ACCESS_TOKEN_SECRET"), exp)
}

// GenerateRefreshToken generates a refresh token.
// Refresh tokens have 1 week of lifepspan.
func GenerateRefreshToken(uid string) (string, error) {
	exp := time.Now().Add(time.Hour * 24 * 7)
	return generateToken(uid, os.Getenv("JWT_REFRESH_TOKEN_TYPE"), os.Getenv("JWT_REFRESH_TOKEN_SECRET"), exp)
}

// VerifyAccessToken verifies the given access token.
func VerifyAccessToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET")), nil
	})
	return parsedToken, err
}

// VerifyRefreshToken verifies the given refresh token.
func VerifyRefreshToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")), nil
	})
	return parsedToken, err
}

// generateToken is a private method that generates and signs a jwt.
// Each type of token have their own secret that should be passed.
func generateToken(uid string, tokenType string, secret string, exp time.Time) (string, error) {
	claims := JwtClaims{
		UserId:    uid,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
