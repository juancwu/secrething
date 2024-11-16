package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			err := next(c)
			latency := time.Now().Sub(start)

			var logger *zerolog.Event
			if err != nil {
				logger = log.Error().Err(err)
				// gotta let the error handler run to get the response data
				c.Error(err)
			} else {
				logger = log.Info()
			}

			protocol := req.Proto
			remoteIP := c.RealIP()
			host := req.Host
			method := req.Method
			uri := req.RequestURI
			uriPath := req.URL.Path
			routePath := c.Path()
			requestID := req.Header.Get(echo.HeaderXRequestID)
			referer := req.Referer()
			userAgent := req.UserAgent()
			status := res.Status
			reqContentType := req.Header.Get(echo.HeaderContentType)
			reqContentLength := req.Header.Get(echo.HeaderContentLength)
			resSize := res.Size
			resContentType := res.Header().Get(echo.HeaderContentType)

			logger.
				Dur("latency", latency).
				Str("protocol", protocol).
				Str("remote_ip", remoteIP).
				Str("host", host).
				Str("method", method).
				Str("uri", uri).
				Str("uri_path", uriPath).
				Str("route_path", routePath).
				Str("request_id", requestID).
				Str("referer", referer).
				Str("user_agent", userAgent).
				Int("status", status).
				Str("request_content_type", reqContentType).
				Str("request_content_length", reqContentLength).
				Int64("response_size", resSize).
				Str("response_content_type", resContentType).
				Send()

			return err
		}
	}
}
