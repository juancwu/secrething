package services

import (
	"fmt"
	"konbini/server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// EmailToken represents the claims in a email jwt
type EmailToken struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

// NewEmailToken creates a new email jwt
func NewEmailToken(emailTokenId string, userId string, salt []byte, exp time.Time) (string, error) {
	cfg, err := config.Global()
	if err != nil {
		return "", err
	}

	claims := EmailToken{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp.UTC()),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ID:        emailTokenId,
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	prefix := cfg.GetEmailTokenKey()
	key := combine(prefix, salt)

	return token.SignedString(key)
}

// ParseEmailToken parses a token into an email token claims without verifying the signature.
func ParseEmailToken(token string) (*EmailToken, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	claims := &EmailToken{}
	_, _, err := parser.ParseUnverified(token, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// VerifyEmailToken verifies the email jwt with the given salt. This salt is gotten from
// the database so it is expected to parse the email jwt unverified beforehand to obtain
// the salt from the database.
func VerifyEmailToken(token string, salt []byte) (*EmailToken, error) {
	cfg, err := config.Global()
	if err != nil {
		return nil, err
	}
	prefix := cfg.GetEmailTokenKey()
	key := combine(prefix, salt)
	claims := &EmailToken{}
	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// combine combines two byte arrays into one
func combine(a []byte, b []byte) []byte {
	combined := make([]byte, len(a)+len(b))
	copy(combined, a)
	copy(combined[len(a):], b)
	return combined
}
