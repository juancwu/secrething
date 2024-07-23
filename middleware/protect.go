package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/juancwu/konbini/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	// JWT_CLAIMS is a key to get the claims parsed by the Protect middleware in echo.Context.
	JWT_CLAIMS = "JWT_CLAIMS"
)

var ErrNoJwtClaims error = errors.New("JWT claims not found in echo context")

// Protect is a middleware that is used to protect a route by authorizing clients with a Bearer token.
//
// This middleware will attach the claims into the echo.Context using JWT_CLAIMS.
func Protect() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId := c.Request().Header.Get(echo.HeaderXRequestID)
			path := c.Request().URL.Path
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				log.Error().Str(echo.HeaderXRequestID, requestId).Str("path", path).Msg("Missing authorization header to access route.")
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Error().Str(echo.HeaderXRequestID, requestId).Str("path", path).Msg("Invalid authorization header.")
				return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid authorization header. Only bearer token is supported."})
			}
			// verify jwt
			claims, err := jwt.VerifyAccessToken(parts[1])
			if err != nil {
				log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Str("path", path).Msg("Failed to validate access token.")
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
			}
			c.Set(JWT_CLAIMS, claims)
			return next(c)
		}
	}
}

// Tries to get the jwt claims set by the Protect() middleware in the echo context.
func GetJwtClaimsFromContext(c echo.Context) (*jwt.JwtClaims, error) {
	claims, ok := c.Get(JWT_CLAIMS).(*jwt.JwtClaims)
	if !ok {
		return nil, ErrNoJwtClaims
	}
	return claims, nil
}
