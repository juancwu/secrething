package middleware

import (
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// RequestID serves as a way to generate a random id for a request
// which serves as a way to traceback the request. Also the RequestID middleware
// from echo does not work so...
func RequestID(l int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			randomId, err := gonanoid.New(l)
			if err != nil {
				return err
			}
			c.Request().Header.Add(echo.HeaderXRequestID, randomId)
			return next(c)
		}
	}
}
