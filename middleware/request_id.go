package middleware

import (
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

func UseRequestId() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			requestId, err := gonanoid.New(16)
			if err != nil {
				logger.Error("Failed to generate request id", zap.Error(err))
			} else {
				c.Set(echo.HeaderXRequestID, requestId)
			}
			return next(c)
		}
	}
}
