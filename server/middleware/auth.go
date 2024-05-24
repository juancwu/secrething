package middleware

import (
	"net/http"
	"strings"

	"github.com/juancwu/konbini/server/service"
	"github.com/labstack/echo/v4"
)

func JwtAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or malformed Bearer token")
		}

		token, err := service.VerifyToken(parts[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if claims, ok := token.Claims.(*service.JwtCustomClaims); ok && token.Valid {
			c.Set("token", token)
			c.Set("claims", claims)
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}
}
