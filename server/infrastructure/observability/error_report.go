package observability

import (
	stderrors "errors"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/juancwu/konbini/server/infrastructure/errors"
	"github.com/labstack/echo/v4"
)

// ReportError reports an error to Sentry with appropriate context and categorization
func ReportError(err error, ctx echo.Context) {
	var appErr errors.AppError
	var echoHTTPError *echo.HTTPError
	var message string

	hub := sentryecho.GetHubFromContext(ctx)
	if hub == nil {
		// Fallback to current hub if context doesn't have one
		hub = sentry.CurrentHub()
	}

	scope := hub.Scope()

	// Check if it's our AppError type
	if stderrors.As(err, &appErr) {
		// Set error information
		scope.SetTag("error_type", string(appErr.Type))
		scope.SetContext("Error Context", map[string]interface{}{
			"public_message":  appErr.PublicMessage,
			"private_message": appErr.PrivateMessage,
			"status_code":     appErr.Code,
		})

		if len(appErr.Errors) > 0 {
			scope.SetContext("Error Details", map[string]interface{}{
				"details": appErr.Errors,
			})
		}

		// Set event level based on error type
		switch appErr.Type {
		case errors.ErrorTypeValidation, errors.ErrorTypeNotFound:
			scope.SetLevel(sentry.LevelInfo)
		case errors.ErrorTypeAuthorization, errors.ErrorTypeForbidden, errors.ErrorTypeRateLimit:
			scope.SetLevel(sentry.LevelWarning)
		case errors.ErrorTypeDatabase, errors.ErrorTypeInternal:
			scope.SetLevel(sentry.LevelError)
		default:
			scope.SetLevel(sentry.LevelError)
		}

		// Use original error if available
		if appErr.InternalError != nil {
			hub.CaptureException(appErr.InternalError)
		} else {
			hub.CaptureException(appErr)
		}
	} else if stderrors.As(err, &echoHTTPError) {
		// Handle Echo's built-in HTTP errors
		statusCode := echoHTTPError.Code

		// Try to get the actual error message
		switch msg := echoHTTPError.Message.(type) {
		case string:
			message = msg
		case error:
			message = msg.Error()
		default:
			message = fmt.Sprintf("%v", echoHTTPError.Message)
		}

		// Configure the scope with Echo HTTP error details
		scope.SetTag("status_code", fmt.Sprintf("%d", statusCode))
		scope.SetContext("HTTP Error", map[string]interface{}{
			"code":    statusCode,
			"message": message,
		})

		// Set level based on HTTP status code
		if statusCode >= 400 && statusCode < 500 {
			scope.SetLevel(sentry.LevelWarning)
		} else if statusCode >= 500 {
			scope.SetLevel(sentry.LevelError)
		} else {
			scope.SetLevel(sentry.LevelInfo)
		}

		hub.CaptureException(echoHTTPError)
	} else {
		// For any other error type
		scope.SetTag("error_type", "unhandled")
		scope.SetLevel(sentry.LevelError)
		scope.SetContext("Unhandled Error", map[string]interface{}{
			"status_code": http.StatusInternalServerError,
			"message":     "An unexpected error occurred",
		})

		hub.CaptureException(err)
	}
}

// CaptureMessage sends a message to Sentry with given level and context
func CaptureMessage(message string, ctx echo.Context, level sentry.Level) {
	hub := sentryecho.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	scope := hub.Scope()
	scope.SetLevel(level)

	// Add request ID if available
	if requestID := ctx.Response().Header().Get(echo.HeaderXRequestID); requestID != "" {
		scope.SetTag("request_id", requestID)
	}

	hub.CaptureMessage(message)
}
