package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewValidationError(t *testing.T) {
	errDetails := []string{"Field is required", "Invalid format"}
	err := NewValidationError("Validation failed", errDetails)

	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected error type %s, got %s", ErrorTypeValidation, err.Type)
	}

	if err.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, err.Code)
	}

	if err.PublicMessage != "Validation failed" {
		t.Errorf("Expected public message 'Validation failed', got '%s'", err.PublicMessage)
	}

	if len(err.Errors) != 2 {
		t.Errorf("Expected 2 error details, got %d", len(err.Errors))
	}

	// Test Error() method
	if err.Error() != "Validation failed" {
		t.Errorf("Expected error message 'Validation failed', got '%s'", err.Error())
	}

	// Test with private message
	err.PrivateMessage = "Internal validation error"
	if err.Error() != "Internal validation error" {
		t.Errorf("Expected error message 'Internal validation error', got '%s'", err.Error())
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("User")

	if err.Type != ErrorTypeNotFound {
		t.Errorf("Expected error type %s, got %s", ErrorTypeNotFound, err.Type)
	}

	if err.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, err.Code)
	}

	if err.PublicMessage != "User not found" {
		t.Errorf("Expected public message 'User not found', got '%s'", err.PublicMessage)
	}
}

func TestNewAuthorizationError(t *testing.T) {
	// Test with custom message
	err := NewAuthorizationError("Invalid token")

	if err.Type != ErrorTypeAuthorization {
		t.Errorf("Expected error type %s, got %s", ErrorTypeAuthorization, err.Type)
	}

	if err.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, err.Code)
	}

	if err.PublicMessage != "Invalid token" {
		t.Errorf("Expected public message 'Invalid token', got '%s'", err.PublicMessage)
	}

	// Test with default message
	err = NewAuthorizationError("")
	if err.PublicMessage != "Invalid or expired credentials" {
		t.Errorf("Expected default message 'Invalid or expired credentials', got '%s'", err.PublicMessage)
	}
}

func TestNewForbiddenError(t *testing.T) {
	// Test with custom message
	err := NewForbiddenError("Admin access required")

	if err.Type != ErrorTypeForbidden {
		t.Errorf("Expected error type %s, got %s", ErrorTypeForbidden, err.Type)
	}

	if err.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, err.Code)
	}

	if err.PublicMessage != "Admin access required" {
		t.Errorf("Expected public message 'Admin access required', got '%s'", err.PublicMessage)
	}

	// Test with default message
	err = NewForbiddenError("")
	if err.PublicMessage != "You do not have permission to perform this action" {
		t.Errorf("Expected default message 'You do not have permission to perform this action', got '%s'", err.PublicMessage)
	}
}

func TestNewInternalError(t *testing.T) {
	originalErr := errors.New("database connection failed")

	// Test with custom message
	err := NewInternalError(originalErr, "System error")

	if err.Type != ErrorTypeInternal {
		t.Errorf("Expected error type %s, got %s", ErrorTypeInternal, err.Type)
	}

	if err.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, err.Code)
	}

	if err.PublicMessage != "System error" {
		t.Errorf("Expected public message 'System error', got '%s'", err.PublicMessage)
	}

	if err.PrivateMessage != "database connection failed" {
		t.Errorf("Expected private message 'database connection failed', got '%s'", err.PrivateMessage)
	}

	if err.InternalError != originalErr {
		t.Errorf("Expected internal error to be the original error")
	}

	// Test with default message
	err = NewInternalError(originalErr, "")
	if err.PublicMessage != "An unexpected error occurred" {
		t.Errorf("Expected default message 'An unexpected error occurred', got '%s'", err.PublicMessage)
	}
}

func TestNewDatabaseError(t *testing.T) {
	originalErr := errors.New("SQL error: table not found")

	// Test with custom message
	err := NewDatabaseError(originalErr, "Database error")

	if err.Type != ErrorTypeDatabase {
		t.Errorf("Expected error type %s, got %s", ErrorTypeDatabase, err.Type)
	}

	if err.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, err.Code)
	}

	if err.PublicMessage != "Database error" {
		t.Errorf("Expected public message 'Database error', got '%s'", err.PublicMessage)
	}

	if err.PrivateMessage != "SQL error: table not found" {
		t.Errorf("Expected private message 'SQL error: table not found', got '%s'", err.PrivateMessage)
	}

	// Test with default message
	err = NewDatabaseError(originalErr, "")
	if err.PublicMessage != "A database error occurred" {
		t.Errorf("Expected default message 'A database error occurred', got '%s'", err.PublicMessage)
	}
}

func TestNewRateLimitError(t *testing.T) {
	// Test with custom message
	err := NewRateLimitError("Too many requests, please try again later")

	if err.Type != ErrorTypeRateLimit {
		t.Errorf("Expected error type %s, got %s", ErrorTypeRateLimit, err.Type)
	}

	if err.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, err.Code)
	}

	if err.PublicMessage != "Too many requests, please try again later" {
		t.Errorf("Expected public message 'Too many requests, please try again later', got '%s'", err.PublicMessage)
	}

	// Test with default message
	err = NewRateLimitError("")
	if err.PublicMessage != "Rate limit exceeded" {
		t.Errorf("Expected default message 'Rate limit exceeded', got '%s'", err.PublicMessage)
	}
}

func TestErrorResponse(t *testing.T) {
	resp := ErrorResponse{
		Code:        http.StatusBadRequest,
		Message:     "Validation failed",
		Errors:      []string{"Field is required"},
		FieldErrors: map[string]interface{}{"name": "Name is required"},
		ReqID:       "req-12345",
	}

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected code %d, got %d", http.StatusBadRequest, resp.Code)
	}

	if resp.Message != "Validation failed" {
		t.Errorf("Expected message 'Validation failed', got '%s'", resp.Message)
	}

	if len(resp.Errors) != 1 || resp.Errors[0] != "Field is required" {
		t.Errorf("Expected errors to contain 'Field is required'")
	}

	if resp.FieldErrors["name"] != "Name is required" {
		t.Errorf("Expected field error for 'name' to be 'Name is required'")
	}

	if resp.ReqID != "req-12345" {
		t.Errorf("Expected req_id 'req-12345', got '%s'", resp.ReqID)
	}
}
