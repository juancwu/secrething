package cookie

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// SetRefreshToken sets an http-only cookie with name "refresh_token"
func SetRefreshToken(c echo.Context, token string) {
	cookie := getRefreshTokenCookieConfig()
	cookie.Value = token
	c.SetCookie(cookie)
}

// UnsetRefreshToken unsets an http-only cookie with name "refresh_token"
func UnsetRefreshToken(c echo.Context) {
	cookie := getRefreshTokenCookieConfig()
	cookie.Value = ""
	c.SetCookie(cookie)
}

// getRefreshTokenCookieConfig gets the default cookie configuration for a refresh token
func getRefreshTokenCookieConfig() *http.Cookie {
	exp := time.Now().UTC().Add(24 * 7 * time.Hour)
	cookie := &http.Cookie{
		Name: "refresh_token",
		// this is a static path, that it should only be allowed in
		Path: "/",
		// TODO: use configuration to determine the domain (CORS)
		Domain: "localhost",
		// TODO: use configuration to determine the secure field (CORS)
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Expires:  exp,
		// Max-Age set to 7 days in seconds
		MaxAge: 7 * 24 * 60 * 60,
	}

	// TODO: Check for is development environment before overrides
	cookie.SameSite = http.SameSiteLaxMode
	// only for dev!!!
	cookie.Secure = false
	cookie.Domain = "localhost"

	return cookie
}
