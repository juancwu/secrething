package middleware

import (
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
)

// SentryHubMiddleware creates a new hub for each request and stores it in the context
// This middleware ensures every request has its own Sentry hub with proper context
func SentryHubMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Clone the current hub to create a request-specific hub
			hub := sentryecho.GetHubFromContext(c)

			// Set basic request information
			hub.Scope().SetRequest(c.Request())

			// Set the request ID as a tag
			if requestID := c.Response().Header().Get(echo.HeaderXRequestID); requestID != "" {
				hub.Scope().SetTag("request_id", requestID)
			}

			// Add remote IP as a tag
			hub.Scope().SetTag("remote_ip", c.RealIP())

			// Add user-agent as a tag
			if userAgent := c.Request().UserAgent(); userAgent != "" {
				hub.Scope().SetTag("user_agent", userAgent)
			}

			// Add route as a tag if available
			if route := c.Path(); route != "" {
				hub.Scope().SetTag("route", route)
			}

			// Add method as a tag
			hub.Scope().SetTag("method", c.Request().Method)

			return next(c)
		}
	}
}
