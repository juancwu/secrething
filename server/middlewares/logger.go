package middlewares

import (
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const CONTEXT_LOGGER_KEY = "middlewares_logger"

type LoggerConfig struct {
	Logger  zerolog.Logger
	Exclude []string
}

func LoggerWithConfig(cfg LoggerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if shouldSkip(c.Path(), cfg.Exclude) {
				return next(c)
			}

			start := time.Now()
			req := c.Request()
			res := c.Response()

			reqLogger := cfg.Logger.With().
				Str("path", req.URL.Path).
				Str("method", req.Method).
				Str("remote_ip", c.RealIP()).
				Str("request_id", req.Header.Get(echo.HeaderXRequestID)).
				Logger()

			c.Set(CONTEXT_LOGGER_KEY, &reqLogger)

			reqLogger.Debug().
				Str("user_agent", req.UserAgent()).
				Interface("headers", req.Header).
				Msg("Incoming request")

			err := next(c)

			duration := time.Since(start)

			logEvent := reqLogger.Info()
			if err != nil {
				logEvent = reqLogger.Error()
			}

			logEvent.
				Int("status_code", res.Status).
				Str("duration", duration.String()).
				Int64("size", res.Size).
				Msg("Request completed")

			return err
		}
	}
}

func GetLogger(c echo.Context) *zerolog.Logger {
	if logger, ok := c.Get(CONTEXT_LOGGER_KEY).(*zerolog.Logger); ok {
		return logger
	}
	defaultLogger := zerolog.New(io.Discard)
	return &defaultLogger
}

func shouldSkip(path string, exclude []string) bool {
	for _, p := range exclude {
		if path == p {
			return true
		}
	}
	return false
}
