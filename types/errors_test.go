package types

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

// TestModel is a struct used for testing validation
type TestModel struct {
	Name     string   `json:"name" validate:"required" errormsg:"Name is mandatory"`
	Email    string   `json:"email" validate:"required,email" errormsg:"required=Email is required;email=Please provide a valid email"`
	Age      int      `json:"age" validate:"min=18,max=100" errormsg:"min=Must be at least 18 years old;max=Must not be older than 100 years"`
	Hobbies  []string `json:"hobbies" validate:"required,min=1" errormsg:"__default=At least one hobby is required"`
	Password string   `json:"password" validate:"required,min=8" errormsg:"required|min=Password must be at least 8 characters long"`
}

func TestNewValidationError(t *testing.T) {
	validate := validator.New()
	model := TestModel{
		Name:     "",
		Email:    "invalid-email",
		Age:      15,
		Hobbies:  []string{},
		Password: "123",
	}

	err := validate.Struct(model)
	assert.Error(t, err)

	errResp := NewValidationError(reflect.TypeOf(model), err)
	assert.NotNil(t, errResp)
	assert.Equal(t, http.StatusBadRequest, errResp.Status)
	assert.Equal(t, "Validation failed", errResp.Message)

	validationErrors, ok := errResp.Details.([]ValidationError)
	assert.True(t, ok)

	// Create a map for easier testing of specific field errors
	errorMap := make(map[string]string)
	for _, ve := range validationErrors {
		errorMap[ve.Field] = ve.Message
	}

	// Verify specific error messages
	assert.Equal(t, "Name is mandatory", errorMap["name"])
	assert.Equal(t, "Please provide a valid email", errorMap["email"])
	assert.Equal(t, "Must be at least 18 years old", errorMap["age"])
	assert.Equal(t, "At least one hobby is required", errorMap["hobbies"])
	assert.Equal(t, "Password must be at least 8 characters long", errorMap["password"])
}

func TestParseErrorMsgTag(t *testing.T) {
	tests := []struct {
		name           string
		structType     reflect.Type
		fieldError     validator.FieldError
		expectedOutput string
	}{
		{
			name:           "Nil struct type",
			structType:     nil,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseErrorMsgTag(tt.structType, tt.fieldError)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func TestContainsFieldTag(t *testing.T) {
	tests := []struct {
		name          string
		errorMsgTag   string
		fieldErrorTag string
		expected      bool
	}{
		{
			name:          "Single matching tag",
			errorMsgTag:   "required",
			fieldErrorTag: "required",
			expected:      true,
		},
		{
			name:          "Multiple tags with match",
			errorMsgTag:   "required|min",
			fieldErrorTag: "min",
			expected:      true,
		},
		{
			name:          "Multiple tags without match",
			errorMsgTag:   "required|min",
			fieldErrorTag: "email",
			expected:      false,
		},
		{
			name:          "Empty tags",
			errorMsgTag:   "",
			fieldErrorTag: "required",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsFieldTag(tt.errorMsgTag, tt.fieldErrorTag)
			assert.Equal(t, tt.expected, result)
		})
	}
}
