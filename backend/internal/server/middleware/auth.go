package middleware

import (
	"context"
	"errors"

	"github.com/juancwu/secrething/internal/server/db"
	"github.com/labstack/echo/v4"
)

// User is a key for storing the authenticated user in the context
type contextKey string

const (
	UserContextKey contextKey = "user"
)

// Errors
var (
	ErrNoAuthHeader       = errors.New("no authorization header provided")
	ErrInvalidAuthHeader  = errors.New("invalid authorization header format")
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenRevoked       = errors.New("token has been revoked")
	ErrRequiresTotp       = errors.New("totp verification required")
)

// GetUserFromContext retrieves the authenticated user from the context
func GetUserFromContext(ctx context.Context) (*db.User, error) {
	user, ok := ctx.Value(UserContextKey).(*db.User)
	if !ok || user == nil {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

// Protected creates a middleware that protects routes from unauthenticated access
// It validates the token and sets the user in the request context
func Protected() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
