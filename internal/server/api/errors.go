package api

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Success      bool   `json:"success"`
	AppErrorCode string `json:"code"`
	ResponseCode int    `json:"-"`
	Message      string `json:"message"`
	RequestID    string `json:"request_id"`
	Err          error  `json:"-"`
}

func (e AppError) Error() string {
	return fmt.Sprintf("AppError (%s - %s): %v", e.Message, e.AppErrorCode, e.Err)
}

func NewBadRequestError(message, code, requestID string, err error) AppError {
	return AppError{
		Success:      false,
		AppErrorCode: code,
		ResponseCode: http.StatusBadRequest,
		Message:      message,
		RequestID:    requestID,
		Err:          err,
	}
}

func NewInternalServerError(message string, requestID string, err error) AppError {
	// Internal server errors typically don't need a specific error code displayed to the user
	// Using a generic code for internal errors
	code := "ERR_INTERNAL_5000"
	return AppError{
		Success:      false,
		AppErrorCode: code,
		ResponseCode: http.StatusInternalServerError,
		Message:      message,
		RequestID:    requestID,
		Err:          err,
	}
}

func NewUnauthorizedError(message, code, requestID string, err error) AppError {
	return AppError{
		Success:      false,
		AppErrorCode: code,
		ResponseCode: http.StatusUnauthorized,
		Message:      message,
		RequestID:    requestID,
		Err:          err,
	}
}

func NewForbiddenError(message, code, requestID string, err error) AppError {
	return AppError{
		Success:      false,
		AppErrorCode: code,
		ResponseCode: http.StatusForbidden,
		Message:      message,
		RequestID:    requestID,
		Err:          err,
	}
}

// NewForbidden is an alias for NewForbiddenError to match the usage pattern
func NewForbidden(message, details, code, requestID string, err error) AppError {
	if details != "" {
		message = fmt.Sprintf("%s: %s", message, details)
	}
	return NewForbiddenError(message, code, requestID, err)
}

func NewNotFoundError(message, code, requestID string, err error) AppError {
	return AppError{
		Success:      false,
		AppErrorCode: code,
		ResponseCode: http.StatusNotFound,
		Message:      message,
		RequestID:    requestID,
		Err:          err,
	}
}

// NewNotFound is an alias for NewNotFoundError to match the usage pattern
func NewNotFound(message, details, code, requestID string, err error) AppError {
	if details != "" {
		message = fmt.Sprintf("%s: %s", message, details)
	}
	return NewNotFoundError(message, code, requestID, err)
}

func NewConflictError(message, code, requestID string, err error) AppError {
	return AppError{
		Success:      false,
		AppErrorCode: code,
		ResponseCode: http.StatusConflict,
		Message:      message,
		RequestID:    requestID,
		Err:          err,
	}
}

// NewConflict is an alias for NewConflictError to match the usage pattern
func NewConflict(message, details, code, requestID string, err error) AppError {
	if details != "" {
		message = fmt.Sprintf("%s: %s", message, details)
	}
	return NewConflictError(message, code, requestID, err)
}

func NewTooManyRequestsError(message, code, requestID string, err error) AppError {
	return AppError{
		Success:      false,
		AppErrorCode: code,
		ResponseCode: http.StatusTooManyRequests,
		Message:      message,
		RequestID:    requestID,
		Err:          err,
	}
}

// NewTooManyRequests is an alias for NewTooManyRequestsError to match the usage pattern
func NewTooManyRequests(message, details, code, requestID string, err error) AppError {
	if details != "" {
		message = fmt.Sprintf("%s: %s", message, details)
	}
	return NewTooManyRequestsError(message, code, requestID, err)
}
