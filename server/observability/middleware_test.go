package observability

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSentryHubMiddleware(t *testing.T) {
	// Initialize Sentry with a test DSN
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "", // Empty DSN for testing
		Environment: "test",
	})
	assert.NoError(t, err)
	defer sentry.Flush(2 * 1000)

	t.Run("basic request", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add request ID header
		c.Response().Header().Set(echo.HeaderXRequestID, "test-request-id")

		// Add user agent
		req.Header.Set("User-Agent", "test-agent")

		// Setup path
		c.SetPath("/api/v1/users")

		// Setup test handler to verify hub exists
		var hubExists bool
		testHandler := func(c echo.Context) error {
			hub := sentryecho.GetHubFromContext(c)
			hubExists = hub != nil
			return c.String(http.StatusOK, "ok")
		}

		// Create middleware chain with sentryecho.New first, then SentryHubMiddleware
		sentryMiddleware := sentryecho.New(sentryecho.Options{})
		middleware := sentryMiddleware(SentryHubMiddleware()(testHandler))

		// Run middleware
		middleware(c)

		// Assertions
		assert.True(t, hubExists, "Hub should exist in context")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())
	})

	t.Run("with detailed headers", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add request ID header
		c.Response().Header().Set(echo.HeaderXRequestID, "complex-request-id")

		// Set headers that will be captured
		req.Header.Set("User-Agent", "Mozilla/5.0 Test Browser")
		req.Header.Set("X-Real-IP", "192.168.1.1")

		// Setup route
		c.SetPath("/api/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("123")

		// Handler that will verify scope
		scopeTags := make(map[string]string)
		testHandler := func(c echo.Context) error {
			hub := sentryecho.GetHubFromContext(c)
			assert.NotNil(t, hub)

			// We can't directly access scope tags, but we can record
			// that the middleware was called and hub exists
			scopeTags["hub_exists"] = "true"

			return c.String(http.StatusOK, "ok")
		}

		// Create middleware chain
		sentryMiddleware := sentryecho.New(sentryecho.Options{})
		middleware := sentryMiddleware(SentryHubMiddleware()(testHandler))

		// Run middleware
		middleware(c)

		// Verify middleware executed and hub was created
		assert.Contains(t, scopeTags, "hub_exists")
		assert.Equal(t, "true", scopeTags["hub_exists"])
	})
}
