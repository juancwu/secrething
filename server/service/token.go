package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/juancwu/konbini/server/env"
)

type JwtCustomClaims struct {
	UserId    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

const (
	ACCESS_TOKEN  = "access_token"
	REFRESH_TOKEN = "refresh_token"
)

func GenerateToken(userId string, tokType string, exp time.Time) (string, error) {
	claims := JwtCustomClaims{
		userId,
		tokType,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    env.Values().JWT_ISSUER,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(env.Values().JWT_SECRET))
}

func GenerateAccessToken(userId string) (string, error) {
	exp := time.Now().Add(time.Hour * 1)
	return GenerateToken(userId, ACCESS_TOKEN, exp)
}

func GenerateRefreshToken(userId string) (string, error) {
	exp := time.Now().Add(time.Hour * 24 * 7) // expires in 1 week
	return GenerateToken(userId, REFRESH_TOKEN, exp)
}

func VerifyToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(env.Values().JWT_SECRET), nil
	})
	return parsedToken, err
}
