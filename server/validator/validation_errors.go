package validator

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
	Value   any    `json:"value,omitempty"`
}

// ValidationErrors is a slice of ValidationError
type ValidationErrors []ValidationError

// Error implements the error interface
func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("validation failed: ")
	for i, err := range v {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return sb.String()
}

// AsMap converts ValidationErrors to a map for consistent error formatting
func (v ValidationErrors) AsMap() map[string]interface{} {
	return FormatValidationErrors(v)
}
