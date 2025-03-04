package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// LoggerMiddleware returns a middleware that logs HTTP requests
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			// Execute the request handler chain
			err := next(c)

			// Calculate request duration
			duration := time.Since(start)

			// Get request ID
			requestID := req.Header.Get(echo.HeaderXRequestID)

			// Check for Cloudflare headers
			cfCountry := req.Header.Get("CF-IPCountry")
			cfConnectingIPv4 := req.Header.Get("CF-Connecting-IP")   // ipv4 of client forwarded from cloudflare
			cfConnectingIPv6 := req.Header.Get("CF-Connecting-IPv6") // ipv6 of client forwarded from cloudflare

			// Check for NGINX headers
			realIP := req.Header.Get("X-Real-IP")
			forwardedFor := req.Header.Get("X-Forwarded-For")
			forwardedProto := req.Header.Get("X-Forwarded-Proto")

			referer := req.Referer()

			// Build the log entry
			logEvent := log.Info().
				Str("request_id", requestID).
				Str("remote_ip", c.RealIP()).
				Str("host", req.Host).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Int64("size", res.Size).
				Dur("duration", duration).
				// Add Cloudflare headers if present
				Str("cf_country", cfCountry).
				Str("cf_connecting_ip", cfConnectingIPv4).
				Str("cf_connecting_ipv6", cfConnectingIPv6).
				// Add NGINX headers if present
				Str("real_ip", realIP).
				Str("forwarded_for", forwardedFor).
				Str("forwarded_proto", forwardedProto).
				// Add user agent
				Str("user_agent", req.UserAgent()).
				// Add referer if present
				Str("referer", referer)

			// Log the entry
			if err != nil {
				// Don't duplicate error logging, as it will be handled by the error handler
				logEvent.Err(err).Msg("Request completed with error")
			} else {
				// Success case
				logEvent.Msg("Request completed successfully")
			}

			return err
		}
	}
}
