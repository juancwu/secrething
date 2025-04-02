package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/services/auth"
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
	tokenService := auth.NewTokenService()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrNoAuthHeader.Error())
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			}

			token := parts[1]

			// Validate the token
			payload, err := tokenService.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidToken.Error())
			}

			// Check if it's a temporary token requiring TOTP verification
			if payload.RequiresTotp {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrRequiresTotp.Error())
			}

			// Get the user from the database
			q, err := db.Query()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to access database")
			}

			user, err := q.GetUserByID(c.Request().Context(), payload.UserID)
			if err != nil {
				if db.IsNoRows(err) {
					return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user data")
			}

			// Store the user in the context for handlers to access
			c.Set(string(UserContextKey), &user)

			return next(c)
		}
	}
}
