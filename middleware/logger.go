package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Logger is a middleware that logs the request and response basic information.
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			path := c.Request().URL.Path
			requestId := c.Request().Header.Get(echo.HeaderXRequestID)
			realIp := c.Request().Header.Get(echo.HeaderXRealIP)
			contentType := c.Request().Header.Get(echo.HeaderContentType)
			contentLength := c.Request().Header.Get(echo.HeaderContentLength)
			// before processing the request
			log.Info().
				Str("method", method).
				Str("path", path).
				Str(echo.HeaderXRequestID, requestId).
				Str(echo.HeaderXRealIP, realIp).
				Str(echo.HeaderContentType, contentType).
				Str(echo.HeaderContentLength, contentLength).
				Send()

				// process the request
			err := next(c)
			// after processing the request
			if err != nil {
				log.Error().
					Err(err).
					Str(echo.HeaderXRequestID, requestId).
					Str(echo.HeaderContentType, c.Response().Header().Get(echo.HeaderContentType)).
					Str(echo.HeaderContentLength, c.Response().Header().Get(echo.HeaderContentLength)).
					Send()
			}
			log.Info().
				Str(echo.HeaderXRequestID, requestId).
				Str(echo.HeaderContentType, c.Response().Header().Get(echo.HeaderContentType)).
				Str(echo.HeaderContentLength, c.Response().Header().Get(echo.HeaderContentLength)).
				Send()
			return err
		}
	}
}
