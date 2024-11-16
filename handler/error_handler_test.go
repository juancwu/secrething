package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"konbini/types"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name           string
		setupError     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "HTTP Error with ErrorResponse",
			setupError: echo.NewHTTPError(http.StatusBadRequest, &types.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid request",
			}),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"invalid request"}`,
		},
		{
			name:           "HTTP Error with string message",
			setupError:     echo.NewHTTPError(http.StatusNotFound, "resource not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"resource not found"}`,
		},
		{
			name:           "Generic error",
			setupError:     errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":500,"message":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// execute
			ErrorHandler(tt.setupError, c)

			// assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestErrorHandler_AlreadyCommitted(t *testing.T) {
	// setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// manually commit the response
	c.NoContent(http.StatusOK)

	// execute
	ErrorHandler(echo.NewHTTPError(http.StatusBadRequest, "test error"), c)

	// assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestErrorHandler_NilMessage(t *testing.T) {
	// setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// execute with nil message
	ErrorHandler(echo.NewHTTPError(http.StatusBadRequest, nil), c)

	// assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), http.StatusText(http.StatusBadRequest))
}
