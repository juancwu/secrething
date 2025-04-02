package auth

import (
	"context"
	"net/http"

	"github.com/juancwu/secrething/internal/server/db"
	handlerErrors "github.com/juancwu/secrething/internal/server/handlers/errors"
	authService "github.com/juancwu/secrething/internal/server/services/auth"
	"github.com/labstack/echo/v4"
)

// generateAndReturnTokens is a helper function that generates and returns tokens
// For web clients, it sets the refresh token as an HTTP-only cookie
// For CLI clients, it returns both tokens in the response body
func generateAndReturnTokens(c echo.Context, ctx context.Context, user *db.User, clientType authService.ClientType, requestID string) error {
	tokenService := authService.NewTokenService()

	// Generate token pair
	tokenPair, err := tokenService.GenerateTokenPair(ctx, user.UserID, clientType)
	if err != nil {
		return handlerErrors.NewInternalServerError("Failed to generate tokens", requestID, err)
	}

	// For web clients, set refresh token as HTTP-only cookie
	if clientType == authService.ClientTypeWeb {
		// Add the cookie settings based on context - HTTP or HTTPS
		secure := c.Scheme() == "https"
		sameSite := http.SameSiteLaxMode
		if secure {
			sameSite = http.SameSiteStrictMode
		}

		// Create the cookie
		cookie := new(http.Cookie)
		cookie.Name = RefreshTokenCookieName
		cookie.Value = tokenPair.RefreshToken
		cookie.Path = RefreshTokenCookiePath
		cookie.MaxAge = CookieMaxAge
		cookie.HttpOnly = true
		cookie.Secure = secure
		cookie.SameSite = sameSite

		// Set the cookie
		c.SetCookie(cookie)

		// Clear refresh token from response - web clients don't need it in JSON
		tokenPair.RefreshToken = ""
	}

	// Return the token response
	return c.JSON(http.StatusOK, loginResponse{
		UserID:       user.UserID,
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken, // Will be empty for web clients
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}
