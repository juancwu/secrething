package middleware

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/juancwu/secrething/internal/server/api"
	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/services"
)

const (
	// User is a key for storing the authenticated user in the context
	UserContextKey string = "user"
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
func GetUserFromContext(ctx echo.Context) (*db.User, error) {
	user, ok := ctx.Get(UserContextKey).(*db.User)
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
			// Extract the token from the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return api.NewUnauthorizedError("No authorization header provided", "ERR_AUTH_REQUIRED_4014", "", ErrNoAuthHeader)
			}

			// Check if the header has the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return api.NewUnauthorizedError("Invalid authorization header format", "ERR_AUTH_REQUIRED_4014", "", ErrInvalidAuthHeader)
			}

			accessToken := parts[1]

			// Verify the token
			ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
			defer cancel()

			tokenService := services.NewTokenService()
			payload, err := tokenService.VerifyToken(ctx, accessToken, services.StdPackage)
			if err != nil {
				switch err := err.(type) {
				case services.TokenServiceError:
					switch err.Type {
					case services.TokenServiceErrExpired:
						return api.NewUnauthorizedError("Token has expired", "ERR_AUTH_EXPIRED_TOKEN_4015", "", err)
					case services.TokenServiceErrInvalid, services.TokenServiceErrDecryption:
						return api.NewUnauthorizedError("Invalid token", "ERR_AUTH_INVALID_TOKEN_4016", "", err)
					}
				}
				return api.NewInternalServerError("Failed to validate token", "", err)
			}

			// Ensure it's an access token
			if payload.TokenType != services.TokenTypeAccess {
				return api.NewUnauthorizedError("Invalid token type", "ERR_AUTH_INVALID_TOKEN_TYPE_4017", "", errors.New("not an access token"))
			}

			// Get the user from database
			q, err := db.Query()
			if err != nil {
				return api.NewInternalServerError("Database error", "", err)
			}

			user, err := q.GetUserByID(ctx, payload.UserID)
			if err != nil {
				if err == sql.ErrNoRows {
					return api.NewUnauthorizedError("User not found", "ERR_AUTH_USER_NOT_FOUND_4011", "", err)
				}
				return api.NewInternalServerError("Failed to get user", "", err)
			}

			// Store user in context
			c.Set(UserContextKey, &user)

			return next(c)
		}
	}
}
