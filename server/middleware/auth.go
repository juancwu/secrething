package middleware

import (
	"net/http"
	"strings"

	"github.com/juancwu/konbini/server/service"
	"github.com/juancwu/konbini/server/utils"
	"github.com/labstack/echo/v4"
)

func BearerAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearerHeader := c.Request().Header.Get("Authorization")
		parts := strings.Split(bearerHeader, " ")
		if len(parts) != 2 {
			utils.Logger().Error("Invalid 'Authorization' header.")
			return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		token := parts[1]
		parsedToken, err := service.VerifyToken(token)
		if err != nil {
			utils.Logger().Errorf("Error verifying token: %v\n", err)
			return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		if !parsedToken.Valid {
			utils.Logger().Error("Invalid token")
			return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		c.Set("token", token)
		c.Set("claims", parsedToken.Claims)

		return next(c)
	}
}
