package middleware

import (
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"
)

const (
	// REQUEST_ID_LENGTH is the length of each request id.
	REQUEST_ID_LENGTH = 12
)

// RequestId is a middleware that generates a random request id of length REQUEST_ID_LENGTH and attaches it to the request header.
func RequestId() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId, err := gonanoid.New(REQUEST_ID_LENGTH)
			if err != nil {
				log.Error().Err(err).Msg("Failed to generate random request id.")
				return next(c)
			}
			c.Request().Header.Set(echo.HeaderXRequestID, requestId)
			return next(c)
		}
	}
}
