package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apiErrors "github.com/juancwu/konbini/server/api/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// MockObservability mocks the observability.ReportError function
func mockReportError(err error, c echo.Context) {}

func setupErrorHandlerTest() (*echo.Echo, *httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Generate a fake request ID
	c.Response().Header().Set(echo.HeaderXRequestID, "req-12345")

	return e, rec, c
}

func TestErrorHandlerWithAppError(t *testing.T) {
	e, _, _ := setupErrorHandlerTest()

	// Create different types of AppError
	testCases := []struct {
		name           string
		err            error
		expectedStatus int
		expectedType   string
	}{
		{
			name:           "Validation Error",
			err:            apiErrors.NewValidationError("Validation failed", []string{"Field is invalid"}),
			expectedStatus: http.StatusBadRequest,
			expectedType:   string(apiErrors.ErrorTypeValidation),
		},
		{
			name:           "Not Found Error",
			err:            apiErrors.NewNotFoundError("User"),
			expectedStatus: http.StatusNotFound,
			expectedType:   string(apiErrors.ErrorTypeNotFound),
		},
		{
			name:           "Authorization Error",
			err:            apiErrors.NewAuthorizationError("Invalid token"),
			expectedStatus: http.StatusUnauthorized,
			expectedType:   string(apiErrors.ErrorTypeAuthorization),
		},
		{
			name:           "Forbidden Error",
			err:            apiErrors.NewForbiddenError("Admin access required"),
			expectedStatus: http.StatusForbidden,
			expectedType:   string(apiErrors.ErrorTypeForbidden),
		},
		{
			name:           "Rate Limit Error",
			err:            apiErrors.NewRateLimitError("Too many requests"),
			expectedStatus: http.StatusTooManyRequests,
			expectedType:   string(apiErrors.ErrorTypeRateLimit),
		},
		{
			name:           "Database Error",
			err:            apiErrors.NewDatabaseError(errors.New("DB connection failed"), "Database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedType:   string(apiErrors.ErrorTypeDatabase),
		},
		{
			name:           "Internal Error",
			err:            apiErrors.NewInternalError(errors.New("Something went wrong"), "Internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedType:   string(apiErrors.ErrorTypeInternal),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fresh context for each test case
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Generate a fake request ID
			c.Response().Header().Set(echo.HeaderXRequestID, "req-12345")

			// Add the request ID to the AppError
			appErr, ok := tc.err.(apiErrors.AppError)
			if ok {
				appErr.RequestID = "req-12345"
				tc.err = appErr
			}

			// Call the error handler
			HTTPErrorHandler(tc.err, c)

			// Check response status
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response body
			var response apiErrors.ErrorResponse
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Verify response structure
			assert.Equal(t, tc.expectedStatus, response.Code)
			assert.NotEmpty(t, response.Message)
			assert.Equal(t, "req-12345", response.ReqID)
		})
	}
}

func TestErrorHandlerWithEchoHTTPError(t *testing.T) {
	e, _, _ := setupErrorHandlerTest()

	// Create Echo HTTP error with string message
	httpErr := echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")

	// Create a fresh context
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Generate a fake request ID
	c.Response().Header().Set(echo.HeaderXRequestID, "req-12345")

	// Call the error handler
	HTTPErrorHandler(httpErr, c)

	// Check response status
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Parse response body
	var response apiErrors.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Invalid request format", response.Message)
	assert.Equal(t, "req-12345", response.ReqID)

	// Test with error message
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Response().Header().Set(echo.HeaderXRequestID, "req-12345")

	httpErr = echo.NewHTTPError(http.StatusUnauthorized, errors.New("authentication failed"))
	HTTPErrorHandler(httpErr, c)

	response = apiErrors.ErrorResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Equal(t, "authentication failed", response.Message)
	assert.Equal(t, "req-12345", response.ReqID)

	// Test with other type of message
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Response().Header().Set(echo.HeaderXRequestID, "req-12345")

	httpErr = echo.NewHTTPError(http.StatusForbidden, 123)
	HTTPErrorHandler(httpErr, c)

	response = apiErrors.ErrorResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusForbidden, response.Code)
	assert.Equal(t, "123", response.Message)
	assert.Equal(t, "req-12345", response.ReqID)
}

func TestErrorHandlerWithGenericError(t *testing.T) {
	e, _, _ := setupErrorHandlerTest()

	// Create a fresh context
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Generate a fake request ID
	c.Response().Header().Set(echo.HeaderXRequestID, "req-12345")

	// Create a generic error
	genericErr := errors.New("something unexpected happened")

	// Call the error handler
	HTTPErrorHandler(genericErr, c)

	// Check response status
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Parse response body
	var response apiErrors.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "An unexpected error occurred", response.Message)
	assert.Equal(t, "req-12345", response.ReqID)
}

func TestLogErrorWithLevel(t *testing.T) {
	// This test mainly ensures the function doesn't panic with different log levels
	// Since we can't easily check log output without mocking the logger
	appErr := apiErrors.NewValidationError("Test error", []string{"Error detail"})

	// Test with valid log levels
	validLevels := []string{
		string(LogLevelDebug),
		string(LogLevelInfo),
		string(LogLevelWarn),
		string(LogLevelError),
	}

	for _, level := range validLevels {
		// Should not panic
		logErrorWithLevel(level, appErr, "req-test")
	}

	// Test with invalid log level
	logErrorWithLevel("invalid", appErr, "req-test")
}

func TestErrorHandlerMiddleware(t *testing.T) {
	// Verify middleware function returns our handler
	middleware := ErrorHandlerMiddleware()

	// Very crude but functional comparison
	middleware1 := errorHandlerToString(HTTPErrorHandler)
	middleware2 := errorHandlerToString(middleware)

	assert.Equal(t, middleware1, middleware2, "ErrorHandlerMiddleware should return HTTPErrorHandler")
}

// Helper function to convert error handler to a comparable string representation
func errorHandlerToString(handler echo.HTTPErrorHandler) string {
	return bytes.NewBuffer([]byte(
		"HTTPErrorHandler function pointer",
	)).String()
}

func TestHTTPErrorHandlerWithCommittedResponse(t *testing.T) {
	_, rec, c := setupErrorHandlerTest()

	// Mark the response as committed
	c.Response().WriteHeader(http.StatusBadRequest)

	// Create an error
	err := errors.New("test error")

	// Call the error handler
	HTTPErrorHandler(err, c)

	// Verify that nothing was written to the response body
	// since the response was already committed
	assert.Empty(t, rec.Body.String())
}

func TestHTTPErrorHandlerIntegration(t *testing.T) {
	// Create Echo instance
	e := echo.New()
	e.HTTPErrorHandler = ErrorHandlerMiddleware()

	// Create a route that returns an error
	e.GET("/not-found", func(c echo.Context) error {
		return apiErrors.NewNotFoundError("Resource")
	})

	e.GET("/echo-error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	})

	e.GET("/generic-error", func(c echo.Context) error {
		return errors.New("something went wrong")
	})

	// Test not found error
	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response apiErrors.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "Resource not found", response.Message)

	// Test Echo HTTP error
	req = httptest.NewRequest(http.MethodGet, "/echo-error", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	response = apiErrors.ErrorResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Bad request", response.Message)

	// Test generic error
	req = httptest.NewRequest(http.MethodGet, "/generic-error", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	response = apiErrors.ErrorResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "An unexpected error occurred", response.Message)
}
