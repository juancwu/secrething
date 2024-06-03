package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}
			defer logger.Sync()
			logger.Info("New request",
				zap.String("Method", c.Request().Method),
				zap.String("Route", c.Request().URL.Path),
				zap.String("Request ID", c.Request().Header.Get(echo.HeaderXRequestID)),
				zap.String("User Agent", c.Request().UserAgent()),
				zap.String("IP", c.RealIP()),
			)
			return next(c)
		}
	}
}
