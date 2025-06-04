package api

import (
	"net/http"
	"strings"

	"github.com/juancwu/secrething/internal/auth"
	"github.com/juancwu/secrething/internal/db"
	"github.com/labstack/echo/v4"
)

// User is the authenticated user context key
type User struct {
	ID    db.UserID `json:"id"`
	Email string    `json:"email"`
}

// AuthMiddleware creates a middleware for authenticating requests
func (api *API) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			// Try to get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				// Check if it's a Bearer token
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}

			// If not found in header, try to get from cookie
			if tokenString == "" {
				cookie, err := c.Cookie("auth_token")
				if err == nil && cookie.Value != "" {
					tokenString = cookie.Value
				}
			}

			// If no token found, return unauthorized
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, apiResponse{
					Code:    http.StatusUnauthorized,
					Message: "Authentication required",
				})
			}

			// Validate token
			claims, err := auth.ValidateToken(tokenString, api.Config.Auth.JWT.Secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, apiResponse{
					Code:    http.StatusUnauthorized,
					Message: "Invalid or expired token",
				})
			}

			// Set user in context
			c.Set("user", User{
				ID:    claims.UserID,
				Email: claims.Email,
			})

			// Continue
			return next(c)
		}
	}
}
