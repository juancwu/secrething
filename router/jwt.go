package router

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// jwtTokenType represents the type of a jwt, which can be of the values JWT_ACCESS_TOKEN or JWT_REFRESH_TOKEN.
type jwtTokenType string

// jwtAuthClaims represents the claims for a signed jwt that is given to a user when logging in.
type jwtAuthClaims struct {
	UserId         string       `json:"user_id"`
	TokenType      jwtTokenType `json:"token_type"`
	ServiceVersion string       `json:"service_version"`
	jwt.RegisteredClaims
}

const (
	// JWT_ACCESS_TOKEN_TYPE is a value used to describe or set the type of a jwt to be an access token. Which is short lived.
	JWT_ACCESS_TOKEN_TYPE jwtTokenType = "konbini_access_token"
	// JWT_ACCESS_TOKEN_TYPE is a value used to describe or set the type of a jwt to be an refresh token. Which is short lived.
	JWT_REFRESH_TOKEN_TYPE jwtTokenType = "konbini_refresh_token"
	// JWT_ACCESS_TOKEN_EXP represents when an access token should expire. It holds the number of nanoseconds for 1 hour. time.Hour * 1
	JWT_ACCESS_TOKEN_EXP int64 = 3600000000000
	// JWT_REFRESH_TOKEN_EXP represents when a refresh token should expire. It holds the number of nanoseconds for 1 week. time.Hour * 24 * 7
	JWT_REFRESH_TOKEN_EXP int64 = 604800000000000
)

// generateToken is a helper function to generate a signed jwt.
// It will not decide for itself when a token should expire and what type of jwt it is.
func generateToken(userId string, tokType jwtTokenType, exp time.Time) (string, error) {
	claims := jwtAuthClaims{
		UserId:         userId,
		TokenType:      tokType,
		ServiceVersion: os.Getenv("VERSION"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// generateAccessToken is a helper function that generates a signed access token in the form of JWT.
// The expiration time of an access token is 1 hour from the time it was created.
func generateAccessToken(userId string) (string, error) {
	exp := time.Now().Add(time.Duration(JWT_ACCESS_TOKEN_EXP))
	return generateToken(userId, JWT_ACCESS_TOKEN_TYPE, exp)
}

// generateRefreshToken is a helper function that generates a signed access token in the form of JWT.
// The expiration time of an refresh token is 1 hour from the time it was created.
func generateRefreshToken(userId string) (string, error) {
	exp := time.Now().Add(time.Duration(JWT_REFRESH_TOKEN_EXP))
	return generateToken(userId, JWT_REFRESH_TOKEN_TYPE, exp)
}

// verifyJWT is a helper function that verifies a jwt
func verifyJWT(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwtAuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return parsedToken, err
}

// useJWT is a helper function that extracts the jwt in the request authorization header and validates the token.
// Explicitly calling this method makes it easier to spot what is actually going on in each route handler.
func useJWT(c echo.Context, acceptedType jwtTokenType) (*jwtAuthClaims, error) {
	authHeaderString := c.Request().Header.Get(echo.HeaderAuthorization)
	parts := strings.Split(authHeaderString, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, errors.New("Invalid authorization header.")
	}
	token, err := verifyJWT(parts[1])
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwtAuthClaims)
	if !ok {
		return nil, errors.New("Invalid jwt casting type.")
	}
	return claims, nil
}
