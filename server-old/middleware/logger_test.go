package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup test handler
	testHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Capture log output
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

	// Add custom headers
	req.Header.Set(echo.HeaderXRequestID, "test-request-id")
	req.Header.Set("CF-IPCountry", "US")
	req.Header.Set("CF-Connecting-IP", "192.0.2.1")
	req.Header.Set("CF-Connecting-IPv6", "2001:db8::1")
	req.Header.Set("X-Real-IP", "192.0.2.2")
	req.Header.Set("X-Forwarded-For", "192.0.2.3")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("Referer", "https://example.com")
	req.Header.Set("User-Agent", "test-agent")

	// Run the middleware
	withLogger := LoggerMiddleware()(testHandler)
	middleware := echomiddleware.RequestID()(withLogger)
	err := middleware(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())

	// Parse log output
	var logOutput map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logOutput)
	t.Log(logOutput)
	assert.NoError(t, err, "Should be able to parse log output as JSON")

	// Check key log fields
	assert.Equal(t, "test-request-id", logOutput["request_id"])
	assert.Equal(t, "GET", logOutput["method"])
	assert.Equal(t, "/test", logOutput["uri"])
	assert.Equal(t, float64(200), logOutput["status"])

	// Check Cloudflare headers
	assert.Equal(t, "US", logOutput["cf_country"])
	assert.Equal(t, "192.0.2.1", logOutput["cf_connecting_ip"])
	assert.Equal(t, "2001:db8::1", logOutput["cf_connecting_ipv6"])

	// Check NGINX headers
	assert.Equal(t, "192.0.2.2", logOutput["real_ip"])
	assert.Equal(t, "192.0.2.3", logOutput["forwarded_for"])
	assert.Equal(t, "https", logOutput["forwarded_proto"])

	// Check other fields
	assert.Equal(t, "test-agent", logOutput["user_agent"])
	assert.Equal(t, "https://example.com", logOutput["referer"])
	assert.Equal(t, "Request completed successfully", logOutput["message"])
}

func TestLoggerMiddlewareWithError(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup test handler with error
	testHandler := func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "test error")
	}

	// Capture log output
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

	// Run the middleware
	middleware := LoggerMiddleware()(testHandler)
	err := middleware(c)

	// Assertions
	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)

	// Parse log output
	var logOutput map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logOutput)
	assert.NoError(t, err, "Should be able to parse log output as JSON")

	// Check error message
	assert.Equal(t, "Request completed with error", logOutput["message"])
	assert.NotNil(t, logOutput["error"])
}

func TestLoggerMiddlewareEmpty(t *testing.T) {
	// Setup with minimal request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Capture log output
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

	// Run middleware
	err := LoggerMiddleware()(func(c echo.Context) error {
		return nil
	})(c)

	// Assertions
	assert.NoError(t, err)

	// Parse log output
	var logOutput map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logOutput)
	assert.NoError(t, err)

	// Even with minimal request, these fields should be present
	assert.Contains(t, logOutput, "remote_ip")
	assert.Contains(t, logOutput, "method")
	assert.Contains(t, logOutput, "uri")
	assert.Contains(t, logOutput, "status")
	assert.Contains(t, logOutput, "duration")
	assert.Equal(t, "Request completed successfully", logOutput["message"])
}
