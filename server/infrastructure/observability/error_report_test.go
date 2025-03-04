package observability

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	apperrors "github.com/juancwu/konbini/server/infrastructure/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestReportError(t *testing.T) {
	// Initialize Sentry with a test DSN
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "", // Empty DSN for testing
		Environment: "test",
	})
	assert.NoError(t, err)
	defer sentry.Flush(2 * 1000)

	// Setup Echo context
	e := echo.New()

	// Test with AppError
	t.Run("with AppError", func(t *testing.T) {
		// Create request and context
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add Sentry hub to context
		sentryMiddleware := sentryecho.New(sentryecho.Options{})
		hubMiddleware := sentryMiddleware(func(c echo.Context) error {
			appErr := apperrors.NewValidationError(
				"Validation failed",
				[]string{"Field 'name' is required"},
			)

			// This is the function we're testing
			ReportError(appErr, c)
			return nil
		})

		// Execute middleware chain
		hubMiddleware(c)

		// We can't easily assert on what was sent to Sentry in tests
		// since the SDK doesn't expose a test client, but we can verify
		// that the code executes without errors
	})

	// Test with Echo HTTPError
	t.Run("with Echo HTTPError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add Sentry hub to context
		sentryMiddleware := sentryecho.New(sentryecho.Options{})
		hubMiddleware := sentryMiddleware(func(c echo.Context) error {
			echoErr := echo.NewHTTPError(http.StatusNotFound, "Resource not found")

			// This is the function we're testing
			ReportError(echoErr, c)
			return nil
		})

		// Execute middleware chain
		hubMiddleware(c)
	})

	// Test with string message Echo HTTPError
	t.Run("with Echo HTTPError containing string message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add Sentry hub to context
		sentryMiddleware := sentryecho.New(sentryecho.Options{})
		hubMiddleware := sentryMiddleware(func(c echo.Context) error {
			echoErr := echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter")

			// This is the function we're testing
			ReportError(echoErr, c)
			return nil
		})

		// Execute middleware chain
		hubMiddleware(c)
	})

	// Test with error message Echo HTTPError
	t.Run("with Echo HTTPError containing error message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add Sentry hub to context
		sentryMiddleware := sentryecho.New(sentryecho.Options{})
		hubMiddleware := sentryMiddleware(func(c echo.Context) error {
			underlyingErr := errors.New("database connection failed")
			echoErr := echo.NewHTTPError(http.StatusInternalServerError, underlyingErr)

			// This is the function we're testing
			ReportError(echoErr, c)
			return nil
		})

		// Execute middleware chain
		hubMiddleware(c)
	})

	// Test with generic error
	t.Run("with generic error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add Sentry hub to context
		sentryMiddleware := sentryecho.New(sentryecho.Options{})
		hubMiddleware := sentryMiddleware(func(c echo.Context) error {
			err := errors.New("something went wrong")

			// This is the function we're testing
			ReportError(err, c)
			return nil
		})

		// Execute middleware chain
		hubMiddleware(c)
	})

	// Test with nil hub (fallback)
	t.Run("with nil hub (fallback to current)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Note: not using sentryMiddleware, so no hub in context
		err := errors.New("fallback error test")

		// Should use CurrentHub() as fallback
		ReportError(err, c)
	})
}

func TestCaptureMessage(t *testing.T) {
	// Initialize Sentry with a test DSN
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "", // Empty DSN for testing
		Environment: "test",
	})
	assert.NoError(t, err)
	defer sentry.Flush(2 * 1000)

	// Setup Echo context
	e := echo.New()

	t.Run("with hub from context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Set request ID
		c.Response().Header().Set(echo.HeaderXRequestID, "test-request-id")

		// Add Sentry hub to context
		sentryMiddleware := sentryecho.New(sentryecho.Options{})
		hubMiddleware := sentryMiddleware(func(c echo.Context) error {
			// This is the function we're testing
			CaptureMessage("Test message", c, sentry.LevelInfo)
			return nil
		})

		// Execute middleware chain
		hubMiddleware(c)
	})

	t.Run("with nil hub (fallback to current)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Set request ID
		c.Response().Header().Set(echo.HeaderXRequestID, "test-request-id")

		// Note: not using sentryMiddleware, so no hub in context
		// Should use CurrentHub() as fallback
		CaptureMessage("Fallback message test", c, sentry.LevelWarning)
	})
}
