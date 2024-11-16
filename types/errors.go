package types

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationError creates a formatted validation error response
func NewValidationError(structType reflect.Type, err error) *ErrorResponse {
	var validationErrors []ValidationError

	if valErrors, ok := err.(validator.ValidationErrors); ok {
		for _, valError := range valErrors {
			validationErrors = append(validationErrors, ValidationError{
				Field:   strings.ToLower(valError.Field()),
				Message: getValidationErrorMsg(structType, valError),
			})
		}
	}

	return &ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "Validation failed",
		Details: validationErrors,
	}
}

// getValidationErrorMsg returns a human-readable validation error message
func getValidationErrorMsg(structType reflect.Type, err validator.FieldError) string {
	msg := parseErrorMsgTag(structType, err)
	if msg != "" {
		return msg
	}
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("Must not be longer than %s characters", err.Param())
	default:
		return fmt.Sprintf("Failed on %s validation", err.Tag())
	}
}

// Parse a struct with the 'errormsg' tag that defines messages for the validator.FieldError.
// Returns the message based on the 'errormsg' format and if the validator.FieldError can be
// matched or not. Otherwise, it returns an empty string.
func parseErrorMsgTag(structType reflect.Type, fieldError validator.FieldError) string {
	var field reflect.StructField
	var found bool

	// early return if no struct type provided.
	// cannot parse error tag
	if structType == nil {
		return ""
	}

	// regex to match all square brackets
	re, err := regexp.Compile(`\[(\d+)\]`)
	if err != nil {
		log.Error().Err(err).Str("struct_name", structType.Name()).Msg("Failed to compiled regex to replace square brackets.")
		return ""
	}
	// remove any square brackets, they come up when there is a field of type slice.
	fields := strings.Split(re.ReplaceAllString(fieldError.StructNamespace(), ""), ".")
	for i, fieldName := range fields {
		// continue since this is the struct name
		if i == 0 {
			continue
		}
		field, found = structType.FieldByName(fieldName)
		if !found {
			return ""
		}
		// get the new struct
		structType = field.Type
		if structType.Kind() == reflect.Slice {
			structType = structType.Elem()
		}
	}

	errormsg := field.Tag.Get("errormsg")

	validationTags := strings.Split(errormsg, ";")
	if len(validationTags) == 1 {
		parts := strings.Split(validationTags[0], "=")
		if len(parts) == 1 {
			// treat as default global message
			return parts[0]
		} else if len(parts) == 2 && (containsFieldTag(parts[0], fieldError.Tag()) || parts[0] == "__default") {
			return parts[1]
		}
	} else if len(validationTags) > 1 {
		for _, tag := range validationTags {
			parts := strings.Split(tag, "=")
			if len(parts) == 2 && (containsFieldTag(parts[0], fieldError.Tag()) || parts[0] == "__default") {
				return parts[1]
			}
		}
	}
	return ""
}

// Checks if the errormsg tag whether or not contains the field error tag.
// This method helps with checking combined errormsg tags.
func containsFieldTag(errorMsgTag, fieldErrorTag string) bool {
	tags := strings.Split(errorMsgTag, "|")
	for _, tag := range tags {
		if tag == fieldErrorTag {
			return true
		}
	}
	return false
}
