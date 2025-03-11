package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error that occurred
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeNotFound represents resource not found errors
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeAuthorization represents authorization errors
	ErrorTypeAuthorization ErrorType = "authorization"
	// ErrorTypeForbidden represents forbidden access errors
	ErrorTypeForbidden ErrorType = "forbidden"
	// ErrorTypeRateLimit represents rate limiting errors
	ErrorTypeRateLimit ErrorType = "rate_limit"
	// ErrorTypeDatabase represents database errors
	ErrorTypeDatabase ErrorType = "database"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "internal"
)

// AppError represents a structured application error
type AppError struct {
	Type           ErrorType // The type of error
	Code           int       // HTTP status code
	RequestID      string    // Request ID for tracking
	PublicMessage  string    // Message that can be shown to the user
	PrivateMessage string    // Internal message with more details (not shown to user)
	Errors         []string  // List of additional error details
	InternalError  error     // The original error if wrapped
}

// Error implements the error interface
func (e AppError) Error() string {
	if e.PrivateMessage != "" {
		return e.PrivateMessage
	}
	return e.PublicMessage
}

// ErrorResponse is the structure sent to the client in error responses
type ErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
	ReqID   string   `json:"request_id,omitempty"`
}

// FieldErrors represents validation errors mapped by field name
type FieldErrors map[string]string

// NewValidationError creates a new validation error
func NewValidationError(publicMessage string, errors []string) AppError {
	return AppError{
		Type:          ErrorTypeValidation,
		Code:          http.StatusBadRequest,
		PublicMessage: publicMessage,
		Errors:        errors,
	}
}

// NewValidationErrorWithFieldErrors creates a new validation error with field-specific error messages
func NewValidationErrorWithFieldErrors(publicMessage string, errors []string, fieldErrors map[string]string) AppError {
	appErr := NewValidationError(publicMessage, errors)
	// Store the field errors in InternalError so they can be accessed by the error handler
	appErr.InternalError = FieldValidationError{FieldErrors: fieldErrors}
	return appErr
}

// FieldValidationError is a custom error type that holds field-specific validation errors
type FieldValidationError struct {
	FieldErrors map[string]string
}

// Error implements the error interface
func (e FieldValidationError) Error() string {
	return "Field validation error"
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string) AppError {
	return AppError{
		Type:          ErrorTypeNotFound,
		Code:          http.StatusNotFound,
		PublicMessage: fmt.Sprintf("%s not found", resource),
	}
}

// NewAuthorizationError creates a new authorization error
func NewAuthorizationError(message string) AppError {
	if message == "" {
		message = "Invalid or expired credentials"
	}
	return AppError{
		Type:          ErrorTypeAuthorization,
		Code:          http.StatusUnauthorized,
		PublicMessage: message,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) AppError {
	if message == "" {
		message = "You do not have permission to perform this action"
	}
	return AppError{
		Type:          ErrorTypeForbidden,
		Code:          http.StatusForbidden,
		PublicMessage: message,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(err error, publicMessage string) AppError {
	if publicMessage == "" {
		publicMessage = "An unexpected error occurred"
	}
	privateMessage := "Internal server error"
	if err != nil {
		privateMessage = err.Error()
	}
	return AppError{
		Type:           ErrorTypeInternal,
		Code:           http.StatusInternalServerError,
		PublicMessage:  publicMessage,
		PrivateMessage: privateMessage,
		InternalError:  err,
	}
}

// NewDatabaseError creates a new database error
func NewDatabaseError(err error, publicMessage string) AppError {
	if publicMessage == "" {
		publicMessage = "A database error occurred"
	}
	privateMessage := "Database error"
	if err != nil {
		privateMessage = err.Error()
	}
	return AppError{
		Type:           ErrorTypeDatabase,
		Code:           http.StatusInternalServerError,
		PublicMessage:  publicMessage,
		PrivateMessage: privateMessage,
		InternalError:  err,
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(message string) AppError {
	if message == "" {
		message = "Rate limit exceeded"
	}
	return AppError{
		Type:          ErrorTypeRateLimit,
		Code:          http.StatusTooManyRequests,
		PublicMessage: message,
	}
}
