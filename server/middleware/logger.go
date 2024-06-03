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
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("request_id", c.Request().Header.Get(echo.HeaderXRequestID)),
				zap.String("user_agent", c.Request().UserAgent()),
				zap.String("ip", c.RealIP()),
			)
			return next(c)
		}
	}
}
