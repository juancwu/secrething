package middleware

import (
	stderrors "errors"
	"fmt"
	"net/http"

	"github.com/juancwu/konbini/server/errors"
	"github.com/juancwu/konbini/server/observability"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// LogLevel defines the acceptable log levels
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// ErrorHandlerMiddleware returns a custom error handler middleware for Echo
func ErrorHandlerMiddleware() echo.HTTPErrorHandler {
	return HTTPErrorHandler
}

// HTTPErrorHandler is an Echo HTTPErrorHandler that properly formats errors
func HTTPErrorHandler(err error, c echo.Context) {
	var appErr errors.AppError
	var echoHTTPError *echo.HTTPError
	var statusCode int
	var message string
	var details []string
	var logLevel string

	observability.ReportError(err, c)

	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// Check if it's our AppError type
	if stderrors.As(err, &appErr) {
		statusCode = appErr.Code
		message = appErr.PublicMessage
		details = appErr.Errors

		// Set request ID if it exists
		if requestID != "" && appErr.RequestID == "" {
			appErr.RequestID = requestID
		}

		// Log with appropriate level based on error type
		switch appErr.Type {
		case errors.ErrorTypeValidation, errors.ErrorTypeNotFound, errors.ErrorTypeAuthorization, errors.ErrorTypeForbidden:
			// These are client errors, log as info
			logLevel = string(LogLevelInfo)
			logErrorWithLevel(logLevel, appErr, requestID)
		case errors.ErrorTypeRateLimit:
			// Rate limiting is worth tracking but not alarming
			logLevel = string(LogLevelWarn)
			logErrorWithLevel(logLevel, appErr, requestID)
		case errors.ErrorTypeDatabase, errors.ErrorTypeInternal:
			// Server errors need more attention
			logLevel = string(LogLevelError)
			logErrorWithLevel(logLevel, appErr, requestID)
		default:
			logLevel = string(LogLevelWarn)
			logErrorWithLevel(logLevel, appErr, requestID)
		}
	} else if stderrors.As(err, &echoHTTPError) {
		// Handle Echo's built-in HTTP errors
		statusCode = echoHTTPError.Code

		// Try to get the actual error message
		switch msg := echoHTTPError.Message.(type) {
		case string:
			message = msg
		case error:
			message = msg.Error()
		default:
			message = fmt.Sprintf("%v", echoHTTPError.Message)
		}

		// Log the error
		logLevel = string(LogLevelWarn)
		log.Warn().
			Str("request_id", requestID).
			Int("status", statusCode).
			Err(err).
			Msg(message)
	} else {
		// For any other error type
		statusCode = http.StatusInternalServerError
		message = "An unexpected error occurred"

		// Log the error
		logLevel = string(LogLevelError)
		log.Error().
			Str("request_id", requestID).
			Err(err).
			Msg("Unhandled error")
	}

	// If the response has already been written, do nothing
	if c.Response().Committed {
		return
	}

	// Prepare the error response
	errorResponse := errors.ErrorResponse{
		Code:    statusCode,
		Message: message,
		Errors:  details,
		ReqID:   requestID,
	}

	// Check if this is a field validation error
	var fieldErrors map[string]interface{}
	if appErr.InternalError != nil {
		// Try to extract field errors if available
		if fieldValidationErr, ok := appErr.InternalError.(errors.FieldValidationError); ok {
			fieldErrors = make(map[string]interface{})
			for field, msg := range fieldValidationErr.FieldErrors {
				fieldErrors[field] = msg
			}
			// Add field errors to the response
			if err := c.JSON(statusCode, map[string]interface{}{
				"code":    statusCode,
				"message": message,
				"errors":  fieldErrors,
				"req_id":  requestID,
			}); err != nil {
				log.Error().
					Str("request_id", requestID).
					Err(err).
					Msg("Failed to send field validation error response")
			}
			return
		}
	}

	// Send standard error response
	if err := c.JSON(statusCode, errorResponse); err != nil {
		log.Error().
			Str("request_id", requestID).
			Err(err).
			Msg("Failed to send error response")
	}
}

// logErrorWithLevel logs an AppError with the appropriate log level
func logErrorWithLevel(level string, appErr errors.AppError, requestID string) {
	// Create a log event with all the context
	logCtx := log.With().
		Str("request_id", requestID).
		Str("error_type", string(appErr.Type)).
		Int("status", appErr.Code).
		Str("public_message", appErr.PublicMessage).
		Str("private_message", appErr.PrivateMessage)

	if len(appErr.Errors) > 0 {
		logCtx = logCtx.Strs("details", appErr.Errors)
	}

	if appErr.InternalError != nil {
		logCtx = logCtx.Err(appErr.InternalError)
	}

	// Create logger from context and log at the appropriate level
	logger := logCtx.Logger()

	// Ensure we only use valid log levels
	validLevel := true

	switch LogLevel(level) {
	case LogLevelDebug:
		logger.Debug().Msg("Application error")
	case LogLevelInfo:
		logger.Info().Msg("Application error")
	case LogLevelWarn:
		logger.Warn().Msg("Application error")
	case LogLevelError:
		logger.Error().Msg("Application error")
	default:
		// Invalid level provided, default to info but log that
		// we received an invalid level to help identify issues
		validLevel = false
		logger.Info().
			Str("invalid_log_level", level).
			Msg("Application error (invalid log level provided)")
	}

	// If we determined this was an invalid level, log a warning about it
	// but only in debug builds or environments
	if !validLevel {
		log.Debug().
			Str("request_id", requestID).
			Str("provided_level", level).
			Msg("Invalid log level specified in error handling")
	}
}
