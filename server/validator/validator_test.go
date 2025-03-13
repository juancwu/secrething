package validator

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestUser is a test struct for validation
type TestUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,gte=18"`
	Password string `json:"password" validate:"required,min=8"`
}

// TestUserWithPasswordValidation is a test struct with the password validation tag
type TestUserWithPasswordValidation struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

// Setup creates an Echo instance and registers the validator
func setupTestValidator() (*echo.Echo, *CustomValidator) {
	e := echo.New()
	v := NewCustomValidator()
	e.Validator = v
	return e, v
}

func TestNewCustomValidator(t *testing.T) {
	v := NewCustomValidator()

	if v.validator == nil {
		t.Error("Expected validator to be initialized")
	}

	if v.translator == nil {
		t.Error("Expected translator to be initialized")
	}
}

func TestSetFieldError(t *testing.T) {
	translator := NewErrorTranslator()
	translator.SetFieldError("name", "required", "Custom name required message")

	// Verify the message was stored
	if translator.fieldErrors["name"]["required"] != "Custom name required message" {
		t.Errorf("Expected custom message to be stored")
	}

	// Test translate function
	msg := translator.Translate("name", "required")
	if msg != "Custom name required message" {
		t.Errorf("Expected to get custom message, got %s", msg)
	}
}

func TestSetDefaultError(t *testing.T) {
	translator := NewErrorTranslator()
	translator.SetDefaultError("required", "This field cannot be empty")

	// Verify the message was stored
	if translator.defaultErrors["required"] != "This field cannot be empty" {
		t.Errorf("Expected default message to be stored")
	}

	// Test translate function for a field without specific message
	msg := translator.Translate("unknown_field", "required")
	if msg != "This field cannot be empty" {
		t.Errorf("Expected to get default message, got %s", msg)
	}
}

func TestSetDefaultMessage(t *testing.T) {
	translator := NewErrorTranslator()
	translator.SetDefaultMessage("Generic error")

	// Verify the message was stored
	if translator.defaultMessage != "Generic error" {
		t.Errorf("Expected default message to be set")
	}

	// Test translate function for an unknown tag
	msg := translator.Translate("field", "unknown_tag")
	if msg != "Generic error" {
		t.Errorf("Expected to get generic message, got %s", msg)
	}
}

func TestValidate(t *testing.T) {
	_, v := setupTestValidator()

	// Test valid case
	validUser := TestUser{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      25,
		Password: "password123",
	}

	err := v.Validate(&validUser)
	if err != nil {
		t.Errorf("Expected no validation errors for valid user, got %v", err)
	}

	// Test invalid case
	invalidUser := TestUser{
		Name:     "",
		Email:    "not-an-email",
		Age:      16,
		Password: "short",
	}

	err = v.Validate(&invalidUser)
	if err == nil {
		t.Error("Expected validation error for invalid user, got nil")
	}

	// Verify the error structure
	validationErrors, ok := err.(ValidationErrors)
	if !ok {
		t.Errorf("Expected ValidationErrors type, got %T", err)
	}

	// Should have 4 validation errors
	if len(validationErrors) != 4 {
		t.Errorf("Expected 4 validation errors, got %d", len(validationErrors))
	}

	// Check if field names are correctly extracted from JSON tags
	fieldNames := make(map[string]bool)
	for _, validationErr := range validationErrors {
		fieldNames[validationErr.Field] = true
	}

	expectedFields := []string{"name", "email", "age", "password"}
	for _, field := range expectedFields {
		if !fieldNames[field] {
			t.Errorf("Expected validation error for field '%s'", field)
		}
	}
}

func TestValidateWithCustomMessages(t *testing.T) {
	_, v := setupTestValidator()

	// Set custom error messages
	v.translator.SetFieldError("name", "required", "Please enter your name")
	v.translator.SetFieldError("email", "email", "Email format is invalid")
	v.translator.SetFieldError("age", "gte", "You must be at least 18 years old")
	v.translator.SetFieldError("password", "min", "Password must be at least 8 characters")

	invalidUser := TestUser{
		Name:     "",
		Email:    "not-an-email",
		Age:      16,
		Password: "short",
	}

	err := v.Validate(&invalidUser)
	validationErrors, _ := err.(ValidationErrors)

	// Find and verify custom messages
	for _, validationErr := range validationErrors {
		switch validationErr.Field {
		case "name":
			if validationErr.Tag == "required" && validationErr.Message != "Please enter your name" {
				t.Errorf("Expected custom message for name.required, got '%s'", validationErr.Message)
			}
		case "email":
			if validationErr.Tag == "email" && validationErr.Message != "Email format is invalid" {
				t.Errorf("Expected custom message for email.email, got '%s'", validationErr.Message)
			}
		case "age":
			if validationErr.Tag == "gte" && validationErr.Message != "You must be at least 18 years old" {
				t.Errorf("Expected custom message for age.gte, got '%s'", validationErr.Message)
			}
		case "password":
			if validationErr.Tag == "min" && validationErr.Message != "Password must be at least 8 characters" {
				t.Errorf("Expected custom message for password.min, got '%s'", validationErr.Message)
			}
		}
	}
}

func TestBindAndValidate(t *testing.T) {
	e, _ := setupTestValidator()

	// Create a test request
	jsonBody := `{"name":"John Doe","email":"john@example.com","age":25,"password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test valid case
	user := new(TestUser)
	err := BindAndValidate(c, user)
	if err != nil {
		t.Errorf("Expected no errors for valid request, got %v", err)
	}

	// Verify bound data
	if user.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", user.Name)
	}

	// Test invalid JSON binding
	invalidJSON := `{"name":}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	user = new(TestUser)
	err = BindAndValidate(c, user)
	if err == nil {
		t.Error("Expected binding error for invalid JSON, got nil")
	}

	// Test validation error
	invalidData := `{"name":"","email":"not-an-email","age":16,"password":"short"}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	user = new(TestUser)
	err = BindAndValidate(c, user)
	if err == nil {
		t.Error("Expected validation error for invalid data, got nil")
	}

	validationErrors, ok := err.(ValidationErrors)
	if !ok {
		t.Errorf("Expected ValidationErrors type, got %T", err)
	}

	if len(validationErrors) != 4 {
		t.Errorf("Expected 4 validation errors, got %d", len(validationErrors))
	}
}

func TestValidationContext(t *testing.T) {
	baseValidator := NewCustomValidator()
	baseValidator.translator.SetDefaultError("required", "Base required message")

	// Create a context with custom messages
	ctx := NewValidationContext(baseValidator)
	ctx.SetFieldError("name", "required", "Context-specific name required message")

	// Verify that the base validator's message is unchanged
	baseMsg := baseValidator.translator.Translate("name", "required")
	if baseMsg != "Base required message" {
		t.Errorf("Base validator should have message 'Base required message', got '%s'", baseMsg)
	}

	// Verify that the context has the custom message
	ctxMsg := ctx.translator.Translate("name", "required")
	if ctxMsg != "Context-specific name required message" {
		t.Errorf("Expected context-specific message, got '%s'", ctxMsg)
	}

	// Test chaining API
	ctx.SetFieldError("email", "email", "Invalid email").
		SetDefaultError("min", "Too short").
		SetDefaultMessage("Generic validation error")

	if ctx.translator.Translate("email", "email") != "Invalid email" {
		t.Error("Chained SetFieldError failed")
	}

	if ctx.translator.Translate("any", "min") != "Too short" {
		t.Error("Chained SetDefaultError failed")
	}

	if ctx.translator.Translate("any", "unknown") != "Generic validation error" {
		t.Error("Chained SetDefaultMessage failed")
	}
}

func TestBindAndValidateWithContext(t *testing.T) {
	e, baseValidator := setupTestValidator()

	// Create a context with custom messages
	ctx := NewValidationContext(baseValidator)
	ctx.SetFieldError("name", "required", "Please provide your full name")
	ctx.SetFieldError("email", "email", "Email address is not valid")

	// Create a test request with invalid data
	invalidData := `{"name":"","email":"not-an-email"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test with context
	user := new(TestUser)
	err := BindAndValidateWithContext(c, user, ctx)

	validationErrors, ok := err.(ValidationErrors)
	if !ok {
		t.Errorf("Expected ValidationErrors type, got %T", err)
	}

	// Find and verify context-specific messages
	for _, validationErr := range validationErrors {
		switch validationErr.Field {
		case "name":
			if validationErr.Tag == "required" {
				assert.Equal(t, "Please provide your full name", validationErr.Message, "Expected context message for name.required")
			}
		case "email":
			if validationErr.Tag == "email" {
				assert.Equal(t, "Email address is not valid", validationErr.Message, "Expected context message for email.email")
			}
		}
	}
}

func TestValidationErrorsAsMap(t *testing.T) {
	errors := ValidationErrors{
		{Field: "name", Message: "Name is required", Tag: "required"},
		{Field: "email", Message: "Invalid email format", Tag: "email"},
	}

	errMap := errors.AsMap()

	if len(errMap) != 2 {
		t.Errorf("Expected map with 2 entries, got %d", len(errMap))
	}

	if errMap["name"] != "Name is required" {
		t.Errorf("Expected message 'Name is required' for field 'name', got '%s'", errMap["name"])
	}

	if errMap["email"] != "Invalid email format" {
		t.Errorf("Expected message 'Invalid email format' for field 'email', got '%s'", errMap["email"])
	}
}

func TestFormatValidationErrors(t *testing.T) {
	errors := ValidationErrors{
		{Field: "name", Message: "Name is required", Tag: "required"},
		{Field: "email", Message: "Invalid email format", Tag: "email"},
		{Field: "age", Message: "Must be at least 18", Tag: "min"},
	}

	fieldErrors := FormatValidationErrors(errors)

	// Check the output structure
	assert.Equal(t, 3, len(fieldErrors), "Expected 3 entries in the field errors map")
	assert.Equal(t, "Name is required", fieldErrors["name"])
	assert.Equal(t, "Invalid email format", fieldErrors["email"])
	assert.Equal(t, "Must be at least 18", fieldErrors["age"])
}

// Define nested struct types for testing
type Address struct {
	Street    string `json:"street" validate:"required"`
	City      string `json:"city" validate:"required"`
	Zip       string `json:"zip" validate:"required"`
	FirstName string `json:"firstName" validate:"required"` // Deliberately same field name as in Profile for testing
}

type Profile struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Bio       string `json:"bio" validate:"max=100"`
}

// Attribute represents an attribute of an item
type Attribute struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// Item represents an item in a list
type Item struct {
	Name       string      `json:"name" validate:"required"`
	Quantity   int         `json:"quantity" validate:"min=1"`
	Price      float64     `json:"price" validate:"required,gt=0"`
	Attributes []Attribute `json:"attributes" validate:"dive"`
	SubItems   []SubItem   `json:"subItems" validate:"dive"`
}

// SubItem represents a nested item within an item
type SubItem struct {
	Name     string `json:"name" validate:"required"`
	Quantity int    `json:"quantity" validate:"min=1"`
}

type UserWithNestedFields struct {
	Username string  `json:"username" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Profile  Profile `json:"profile" validate:"required"`
	Address  Address `json:"address" validate:"required"`
	Items    []Item  `json:"items" validate:"required,dive"`   // Add items slice with nested structures
	Matrix   [][]int `json:"matrix" validate:"dive,dive,gt=0"` // 2D array for testing multiple dive levels
}

func TestNestedValidationFormatting(t *testing.T) {
	// Create an instance with validation errors
	invalidUser := UserWithNestedFields{
		Username: "",             // required error
		Email:    "not-an-email", // email error
		Profile: Profile{
			FirstName: "",                                                                                                                                       // required error
			LastName:  "",                                                                                                                                       // required error
			Bio:       "This is a very long bio that exceeds the maximum length limit of 100 characters. It should trigger a validation error for the max tag.", // max error
		},
		Address: Address{
			Street: "", // required error
			City:   "", // required error
			Zip:    "", // required error
		},
	}

	// Create validator and validate
	validator := NewCustomValidator()
	err := validator.Validate(&invalidUser)
	assert.NotNil(t, err, "Expected validation error")

	// Convert to ValidationErrors
	validationErrors, ok := err.(ValidationErrors)
	assert.True(t, ok, "Expected ValidationErrors type")
	assert.True(t, len(validationErrors) > 0, "Expected at least one validation error")

	// Format the errors
	fieldErrors := FormatValidationErrors(validationErrors)

	// Verify top-level fields
	assert.Contains(t, fieldErrors, "username")
	assert.Contains(t, fieldErrors, "email")

	// Verify nested profile fields
	profileMap, ok := fieldErrors["profile"].(map[string]interface{})
	assert.True(t, ok, "Expected 'profile' to be a nested map")
	assert.Contains(t, profileMap, "firstName")
	assert.Contains(t, profileMap, "lastName")
	assert.Contains(t, profileMap, "bio")

	// Verify nested address fields
	addressMap, ok := fieldErrors["address"].(map[string]interface{})
	assert.True(t, ok, "Expected 'address' to be a nested map")
	assert.Contains(t, addressMap, "street")
	assert.Contains(t, addressMap, "city")
	assert.Contains(t, addressMap, "zip")

	// Output formatted errors for debugging
	jsonBytes, _ := json.MarshalIndent(fieldErrors, "", "  ")
	t.Logf("Formatted validation errors: %s", string(jsonBytes))
}

func TestCustomMessagesForNestedFields(t *testing.T) {
	// Create a validator
	validator := NewCustomValidator()

	// Set custom messages for nested fields by leaf field name
	validator.translator.SetFieldError("firstName", "required", "First name is mandatory")
	validator.translator.SetFieldError("lastName", "required", "Last name is mandatory")
	validator.translator.SetFieldError("street", "required", "Street address cannot be empty")
	validator.translator.SetFieldError("bio", "max", "Bio is too long")

	// Create an instance with validation errors
	invalidUser := UserWithNestedFields{
		Username: "testuser",      // valid
		Email:    "test@test.com", // valid
		Profile: Profile{
			FirstName: "",                                                                                                                                       // required error
			LastName:  "",                                                                                                                                       // required error
			Bio:       "This is a very long bio that exceeds the maximum length limit of 100 characters. It should trigger a validation error for the max tag.", // max error
		},
		Address: Address{
			Street: "",         // required error
			City:   "New York", // valid
			Zip:    "10001",    // valid
		},
	}

	// Validate
	err := validator.Validate(&invalidUser)
	assert.NotNil(t, err, "Expected validation error")

	// Convert to ValidationErrors
	validationErrors, ok := err.(ValidationErrors)
	assert.True(t, ok, "Expected ValidationErrors type")

	// Verify custom messages
	fieldErrors := make(map[string]string)
	for _, err := range validationErrors {
		fieldErrors[err.Field] = err.Message
	}

	// Check for custom messages by full field path
	assert.Equal(t, "First name is mandatory", fieldErrors["profile.firstName"])
	assert.Equal(t, "Last name is mandatory", fieldErrors["profile.lastName"])
	assert.Equal(t, "Bio is too long", fieldErrors["profile.bio"])
	assert.Equal(t, "Street address cannot be empty", fieldErrors["address.street"])

	// Format the errors into a nested structure
	formattedErrors := FormatValidationErrors(validationErrors)

	// Output formatted errors for debugging
	jsonBytes, _ := json.MarshalIndent(formattedErrors, "", "  ")
	t.Logf("Formatted validation errors: %s", string(jsonBytes))

	// Verify that formatted errors maintain the custom messages
	profileMap, ok := formattedErrors["profile"].(map[string]interface{})
	assert.True(t, ok, "Expected 'profile' to be a nested map")
	assert.Equal(t, "First name is mandatory", profileMap["firstName"])
	assert.Equal(t, "Last name is mandatory", profileMap["lastName"])
	assert.Equal(t, "Bio is too long", profileMap["bio"])

	addressMap, ok := formattedErrors["address"].(map[string]interface{})
	assert.True(t, ok, "Expected 'address' to be a nested map")
	assert.Equal(t, "Street address cannot be empty", addressMap["street"])
}

func TestCustomMessagesForNestedFieldsWithContext(t *testing.T) {
	// Create a base validator
	baseValidator := NewCustomValidator()

	// Create a validation context
	ctx := NewValidationContext(baseValidator)

	// Set custom messages for nested fields using the context
	ctx.SetFieldError("firstName", "required", "Please enter your first name")
	ctx.SetFieldError("lastName", "required", "Please enter your last name")
	ctx.SetFieldError("bio", "max", "Your biography is too detailed, please shorten it")
	ctx.SetFieldError("street", "required", "Please provide a street address")

	// Create an instance with validation errors
	invalidUser := UserWithNestedFields{
		Username: "testuser",      // valid
		Email:    "test@test.com", // valid
		Profile: Profile{
			FirstName: "",                                                                                                                                       // required error
			LastName:  "",                                                                                                                                       // required error
			Bio:       "This is a very long bio that exceeds the maximum length limit of 100 characters. It should trigger a validation error for the max tag.", // max error
		},
		Address: Address{
			Street: "",         // required error
			City:   "New York", // valid
			Zip:    "10001",    // valid
		},
	}

	// Validate using the context
	err := ctx.Validate(&invalidUser)
	assert.NotNil(t, err, "Expected validation error")

	// Convert to ValidationErrors
	validationErrors, ok := err.(ValidationErrors)
	assert.True(t, ok, "Expected ValidationErrors type")

	// Verify custom messages
	fieldErrors := make(map[string]string)
	for _, err := range validationErrors {
		fieldErrors[err.Field] = err.Message
	}

	// Check for custom messages by full field path
	assert.Equal(t, "Please enter your first name", fieldErrors["profile.firstName"])
	assert.Equal(t, "Please enter your last name", fieldErrors["profile.lastName"])
	assert.Equal(t, "Your biography is too detailed, please shorten it", fieldErrors["profile.bio"])
	assert.Equal(t, "Please provide a street address", fieldErrors["address.street"])

	// Verify that the base validator is unchanged
	// Create a new instance for testing the base validator
	anotherInvalidUser := UserWithNestedFields{
		Username: "",
		Email:    "",
		Profile: Profile{
			FirstName: "",
			LastName:  "",
		},
		Address: Address{
			Street: "",
			City:   "",
			Zip:    "",
		},
	}

	baseErr := baseValidator.Validate(&anotherInvalidUser)
	baseValidationErrors, _ := baseErr.(ValidationErrors)

	// The base validator should still use default messages
	for _, err := range baseValidationErrors {
		if err.Field == "profile.firstName" && err.Tag == "required" {
			assert.NotEqual(t, "Please enter your first name", err.Message,
				"Base validator message should not be affected by context")
		}
	}

	// Format the errors into a nested structure using the context validation errors
	formattedErrors := FormatValidationErrors(validationErrors)

	// Verify that formatted errors maintain the custom messages
	profileMap, ok := formattedErrors["profile"].(map[string]interface{})
	assert.True(t, ok, "Expected 'profile' to be a nested map")
	assert.Equal(t, "Please enter your first name", profileMap["firstName"])
	assert.Equal(t, "Please enter your last name", profileMap["lastName"])
	assert.Equal(t, "Your biography is too detailed, please shorten it", profileMap["bio"])

	addressMap, ok := formattedErrors["address"].(map[string]interface{})
	assert.True(t, ok, "Expected 'address' to be a nested map")
	assert.Equal(t, "Please provide a street address", addressMap["street"])
}

// TestGlobalErrorHandler was removed as GlobalErrorHandler was moved to the middleware package

func TestFullPathCustomMessages(t *testing.T) {
	// Create a validator
	validator := NewCustomValidator()

	// Set custom messages for fields with the same leaf name but different paths
	validator.translator.SetFieldError("profile.firstName", "required", "Profile first name is required")
	validator.translator.SetFieldError("address.firstName", "required", "Contact first name is required")

	// Create an instance with validation errors in both fields with the same leaf name
	invalidUser := UserWithNestedFields{
		Username: "testuser",      // valid
		Email:    "test@test.com", // valid
		Profile: Profile{
			FirstName: "",      // required error - should get profile-specific message
			LastName:  "Smith", // valid
		},
		Address: Address{
			Street:    "123 Main St", // valid
			City:      "New York",    // valid
			Zip:       "10001",       // valid
			FirstName: "",            // required error - should get address-specific message
		},
		Items: []Item{},
	}

	// Validate
	err := validator.Validate(&invalidUser)
	assert.NotNil(t, err, "Expected validation error")

	// Convert to ValidationErrors
	validationErrors, ok := err.(ValidationErrors)
	assert.True(t, ok, "Expected ValidationErrors type")

	// Verify that each firstName field gets its specific custom message
	fieldErrors := make(map[string]string)
	for _, err := range validationErrors {
		fieldErrors[err.Field] = err.Message
	}

	// Each path should have its own custom message
	assert.Equal(t, "Profile first name is required", fieldErrors["profile.firstName"])
	assert.Equal(t, "Contact first name is required", fieldErrors["address.firstName"])

	// Format the errors into a nested structure
	formattedErrors := FormatValidationErrors(validationErrors)

	// Output formatted errors for debugging
	jsonBytes, _ := json.MarshalIndent(formattedErrors, "", "  ")
	t.Logf("Formatted validation errors: %s", string(jsonBytes))

	// Verify that formatted errors maintain the path-specific custom messages
	profileMap, ok := formattedErrors["profile"].(map[string]interface{})
	assert.True(t, ok, "Expected 'profile' to be a nested map")
	assert.Equal(t, "Profile first name is required", profileMap["firstName"])

	addressMap, ok := formattedErrors["address"].(map[string]interface{})
	assert.True(t, ok, "Expected 'address' to be a nested map")
	assert.Equal(t, "Contact first name is required", addressMap["firstName"])
}

func TestCloneValidator(t *testing.T) {
	original := NewCustomValidator()
	original.translator.SetDefaultError("required", "Original required message")
	original.translator.SetFieldError("name", "required", "Original name required message")

	// Clone the validator
	clone := original.Clone()

	// Modify the clone
	clone.translator.SetDefaultError("required", "Cloned required message")
	clone.translator.SetFieldError("name", "required", "Cloned name required message")

	// Verify original is unchanged
	if original.translator.Translate("any", "required") != "Original required message" {
		t.Error("Original default error message was modified")
	}

	if original.translator.Translate("name", "required") != "Original name required message" {
		t.Error("Original field error message was modified")
	}

	// Verify clone has new messages
	if clone.translator.Translate("any", "required") != "Cloned required message" {
		t.Error("Clone default error message not set correctly")
	}

	if clone.translator.Translate("name", "required") != "Cloned name required message" {
		t.Error("Clone field error message not set correctly")
	}
}

// Integration test with Echo
func TestPasswordValidation(t *testing.T) {
	validator := NewCustomValidator()

	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{
			name:     "Valid password with all requirements",
			password: "Password123!",
			valid:    true,
		},
		{
			name:     "Valid password with minimum requirements",
			password: "Aa1!5678",
			valid:    true,
		},
		{
			name:     "Too short",
			password: "Aa1!",
			valid:    false,
		},
		{
			name:     "Missing uppercase",
			password: "password123!",
			valid:    false,
		},
		{
			name:     "Missing lowercase",
			password: "PASSWORD123!",
			valid:    false,
		},
		{
			name:     "Missing digit",
			password: "Password!",
			valid:    false,
		},
		{
			name:     "Missing special character",
			password: "Password123",
			valid:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user := TestUserWithPasswordValidation{
				Username: "testuser",
				Email:    "test@example.com",
				Password: tc.password,
			}

			err := validator.Validate(&user)

			if tc.valid {
				assert.NoError(t, err, "Expected no validation error")
			} else {
				assert.Error(t, err, "Expected validation error")

				if err != nil {
					validationErrors, ok := err.(ValidationErrors)
					assert.True(t, ok, "Expected ValidationErrors type")

					// Find password validation error
					var passwordError *ValidationError
					for _, validationErr := range validationErrors {
						if validationErr.Field == "password" && validationErr.Tag == "password" {
							passwordError = &validationErr
							break
						}
					}

					assert.NotNil(t, passwordError, "Expected password validation error")
					if passwordError != nil {
						assert.Equal(t, "password", passwordError.Tag)
						assert.Contains(t, passwordError.Message, "Password must be at least 8 characters long")
					}
				}
			}
		})
	}
}

func TestDefaultValidationMessages(t *testing.T) {
	validator := NewCustomValidator()

	// Define the expected default messages
	expectedDefaults := map[string]string{
		"required": "This field is required",
		"email":    "Must be a valid email address",
		"min":      "Value must be greater than or equal to the minimum",
		"max":      "Value must be less than or equal to the maximum",
		"len":      "Must have the exact required length",
		"eq":       "Value must be equal to the required value",
		"ne":       "Value cannot be equal to the specified value",
		"oneof":    "Must be one of the available options",
		"url":      "Must be a valid URL",
		"alpha":    "Must contain only letters",
		"alphanum": "Must contain only letters and numbers",
		"numeric":  "Must be a valid numeric value",
		"uuid":     "Must be a valid UUID",
		"datetime": "Must be a valid date/time",
		"password": "Password must be at least 8 characters long and contain uppercase, lowercase, digit, and at least one special character (!@#$%^&*()-_=+[]{}|;:'\",.<>/?)",
	}

	// Verify all default messages are set correctly
	for tag, expectedMessage := range expectedDefaults {
		actualMessage := validator.translator.Translate("any_field", tag)
		assert.Equal(t, expectedMessage, actualMessage, "Default message for '%s' tag doesn't match expected value", tag)
	}

	// Test validation with default messages
	invalidUser := TestUser{
		Name:     "",
		Email:    "not-an-email",
		Age:      16,
		Password: "short",
	}

	err := validator.Validate(&invalidUser)
	assert.NotNil(t, err, "Expected validation error")

	validationErrors, ok := err.(ValidationErrors)
	assert.True(t, ok, "Expected ValidationErrors type")

	// Check for specific validation messages
	for _, validationErr := range validationErrors {
		switch {
		case validationErr.Field == "name" && validationErr.Tag == "required":
			assert.Equal(t, "This field is required", validationErr.Message)
		case validationErr.Field == "email" && validationErr.Tag == "email":
			assert.Equal(t, "Must be a valid email address", validationErr.Message)
		case validationErr.Field == "password" && validationErr.Tag == "min":
			assert.Equal(t, "Value must be greater than or equal to the minimum", validationErr.Message)
		}
	}
}

func TestIntegrationWithEcho(t *testing.T) {
	// Create Echo instance with validator
	e := echo.New()
	v := NewCustomValidator()
	v.translator.SetFieldError("name", "required", "Name cannot be empty")
	v.translator.SetFieldError("email", "email", "Please provide a valid email")
	e.Validator = v

	// Setup test handler
	e.POST("/users", func(c echo.Context) error {
		user := new(TestUser)
		if err := BindAndValidate(c, user); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success": true,
			"user":    user,
		})
	})

	// Set a test error handler that uses FormatValidationErrors
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var validationErrors ValidationErrors
		var statusCode int = http.StatusInternalServerError
		var message string = "Internal server error"

		if stderrors.As(err, &validationErrors) {
			// Handle validation errors
			statusCode = http.StatusBadRequest
			message = "Validation failed"

			// Use the FormatValidationErrors function
			fieldErrors := FormatValidationErrors(validationErrors)

			c.JSON(statusCode, map[string]interface{}{
				"code":         statusCode,
				"message":      message,
				"field_errors": fieldErrors,
			})
			return
		}

		// Handle other errors
		if httpErr, ok := err.(*echo.HTTPError); ok {
			statusCode = httpErr.Code
			message = fmt.Sprintf("%v", httpErr.Message)
		}

		c.JSON(statusCode, map[string]interface{}{
			"code":    statusCode,
			"message": message,
		})
	}

	// Test with invalid data
	invalidData := `{"name":"","email":"not-an-email"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(invalidData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Equal(t, "Validation failed", response["message"])

	fieldErrors, ok := response["field_errors"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Name cannot be empty", fieldErrors["name"])
	assert.Equal(t, "Please provide a valid email", fieldErrors["email"])

	// Test with valid data
	validData := `{"name":"John Doe","email":"john@example.com","age":25,"password":"password123"}`
	req = httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(validData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Verify success response
	assert.Equal(t, http.StatusCreated, rec.Code)

	response = make(map[string]interface{})
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])

	user, ok := response["user"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "John Doe", user["name"])
	assert.Equal(t, "john@example.com", user["email"])
}

func TestArrayOfStructsValidation(t *testing.T) {
	// Define address struct
	type AddressStruct struct {
		Street       string `json:"street" validate:"required"`
		SecondStreet string `json:"second_street" validate:"required"`
		City         string `json:"city" validate:"required"`
		ZipCode      string `json:"zip_code" validate:"required"`
	}

	// Define profile struct with addresses
	type ProfileStruct struct {
		FirstName string          `json:"first_name" validate:"required"`
		LastName  string          `json:"last_name" validate:"required"`
		Addresses []AddressStruct `json:"addresses" validate:"required,dive"`
	}

	// Define the main test struct
	type TestUser struct {
		Username string        `json:"username" validate:"required"`
		Email    string        `json:"email" validate:"required,email"`
		Profile  ProfileStruct `json:"profile" validate:"required"`
	}

	// Create a validator
	validator := NewCustomValidator()

	// Set custom error messages
	validator.translator.SetFieldError("username", "required", "Username is required")
	validator.translator.SetFieldError("email", "email", "Invalid email format")
	validator.translator.SetFieldError("profile.first_name", "required", "First name is required")
	validator.translator.SetFieldError("profile.addresses", "required", "At least one address is required")
	validator.translator.SetFieldError("street", "required", "Street is required")
	validator.translator.SetFieldError("second_street", "required", "Second street is required")
	validator.translator.SetFieldError("city", "required", "City is required")
	validator.translator.SetFieldError("zip_code", "required", "Zip code is required")

	// Create an invalid test user with multiple validation errors in address structs
	invalidUser := TestUser{
		Username: "",
		Email:    "invalid-email",
		Profile: ProfileStruct{
			FirstName: "",
			LastName:  "Doe",
			Addresses: []AddressStruct{
				{
					Street:       "", // Required error
					SecondStreet: "", // Required error
					City:         "New York",
					ZipCode:      "10001",
				},
				{
					Street:       "456 Other St",
					SecondStreet: "Apt 789",
					City:         "", // Required error
					ZipCode:      "", // Required error
				},
			},
		},
	}

	// Validate
	err := validator.Validate(&invalidUser)
	assert.NotNil(t, err, "Expected validation error")

	// Convert to ValidationErrors
	validationErrors, ok := err.(ValidationErrors)
	assert.True(t, ok, "Expected ValidationErrors type")

	// Log all validation errors for debugging
	for _, err := range validationErrors {
		t.Logf("Field: %s, Tag: %s, Message: %s", err.Field, err.Tag, err.Message)
	}

	// Format the errors
	formattedErrors := FormatValidationErrors(validationErrors)

	// Output formatted errors for debugging
	jsonBytes, _ := json.MarshalIndent(formattedErrors, "", "  ")
	t.Logf("Formatted validation errors: %s", string(jsonBytes))

	// Verify the nested structure with arrays of structs

	// 1. Verify top-level fields
	assert.Equal(t, "Username is required", formattedErrors["username"])
	assert.Equal(t, "Invalid email format", formattedErrors["email"])

	// 2. Verify profile fields
	profile, ok := formattedErrors["profile"].(map[string]interface{})
	assert.True(t, ok, "Expected 'profile' to be a map")
	assert.Equal(t, "First name is required", profile["first_name"])

	// 3. Verify addresses array contains structs with their own validation errors
	addresses, ok := profile["addresses"].([]interface{})
	assert.True(t, ok, "Expected 'addresses' to be an array")
	assert.Equal(t, 2, len(addresses), "Expected 2 addresses in the array")

	// 4. Verify first address validation errors with index information
	firstAddress, ok := addresses[0].(map[string]interface{})
	assert.True(t, ok, "Expected first address to be a map")
	assert.Equal(t, 0, firstAddress["index"])

	firstAddressErrors, ok := firstAddress["field_errors"].(map[string]interface{})
	assert.True(t, ok, "Expected field_errors to be a map")
	assert.Equal(t, "Street is required", firstAddressErrors["street"])
	assert.Equal(t, "Second street is required", firstAddressErrors["second_street"])

	// 5. Verify second address validation errors with index information
	secondAddress, ok := addresses[1].(map[string]interface{})
	assert.True(t, ok, "Expected second address to be a map")
	assert.Equal(t, 1, secondAddress["index"])

	secondAddressErrors, ok := secondAddress["field_errors"].(map[string]interface{})
	assert.True(t, ok, "Expected field_errors to be a map")
	assert.Equal(t, "City is required", secondAddressErrors["city"])
	assert.Equal(t, "Zip code is required", secondAddressErrors["zip_code"])

	// 6. Ensure expected structure matches exactly with index information
	expectedStructure := map[string]interface{}{
		"username": "Username is required",
		"email":    "Invalid email format",
		"profile": map[string]interface{}{
			"first_name": "First name is required",
			"addresses": []interface{}{
				map[string]interface{}{
					"index": 0,
					"field_errors": map[string]interface{}{
						"street":        "Street is required",
						"second_street": "Second street is required",
					},
				},
				map[string]interface{}{
					"index": 1,
					"field_errors": map[string]interface{}{
						"city":     "City is required",
						"zip_code": "Zip code is required",
					},
				},
			},
		},
	}

	// Deep equality check against expected structure
	assert.Equal(t, expectedStructure, formattedErrors, "Formatted errors should match expected structure")
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "Simple field",
			input:    "name",
			expected: "name",
		},
		{
			name:     "Nested field",
			input:    "profile.firstName",
			expected: "profile.firstName",
		},
		{
			name:     "Path with array index",
			input:    "addresses[0].street",
			expected: "addresses[0].street",
		},
		{
			name:     "Path with multiple array indices",
			input:    "data[0][1].value",
			expected: "data[0][1].value",
		},
		{
			name:     "Path with array index and nested fields",
			input:    "user.addresses[0].street",
			expected: "user.addresses[0].street",
		},
		{
			name:     "Path with spaces",
			input:    "  user.name  ",
			expected: "user.name",
		},
		{
			name:     "Path with empty segments",
			input:    "user..name",
			expected: "user.name",
		},
		{
			name:     "Path with map key",
			input:    "metadata[\"key\"].value",
			expected: "metadata[\"key\"].value",
		},
		{
			name:     "Path with array and whitespace",
			input:    "items [ 0 ] . name",
			expected: "items[0].name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePath(tt.input)
			assert.Equal(t, tt.expected, result, "Normalized path should match expected value")
		})
	}
}

func TestExtractLeafName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "Simple field",
			input:    "name",
			expected: "name",
		},
		{
			name:     "Nested field",
			input:    "profile.firstName",
			expected: "firstName",
		},
		{
			name:     "Path with array index at leaf",
			input:    "items[0]",
			expected: "items[0]",
		},
		{
			name:     "Path with array index in middle",
			input:    "orders[0].items",
			expected: "items",
		},
		{
			name:     "Path with array index at leaf after dot",
			input:    "data.points[0]",
			expected: "points[0]",
		},
		{
			name:     "Deep nested path",
			input:    "user.profile.contact.email",
			expected: "email",
		},
		{
			name:     "Deep nested path with array",
			input:    "user.addresses[0].street.name",
			expected: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractLeafName(tt.input)
			assert.Equal(t, tt.expected, result, "Extracted leaf name should match expected value")
		})
	}
}

func TestPathNormalizedValidation(t *testing.T) {
	// Create a validator
	validator := NewCustomValidator()

	// Set custom error messages with different path formats
	validator.translator.SetFieldError("profile.addresses[0].street", "required", "Street is required")

	// Create a test struct with nested fields and arrays
	type Address struct {
		Street string `json:"street" validate:"required"`
		City   string `json:"city" validate:"required"`
	}

	type Profile struct {
		Addresses []Address `json:"addresses" validate:"required,dive"`
	}

	type User struct {
		Profile Profile `json:"profile" validate:"required"`
	}

	// Create an invalid user with missing street
	invalidUser := User{
		Profile: Profile{
			Addresses: []Address{
				{
					Street: "", // Required error - should match our custom message
					City:   "New York",
				},
			},
		},
	}

	// Validate
	err := validator.Validate(&invalidUser)
	assert.NotNil(t, err, "Expected validation error")

	// Convert to ValidationErrors
	validationErrors, ok := err.(ValidationErrors)
	assert.True(t, ok, "Expected ValidationErrors type")

	// Find the street validation error
	var streetError *ValidationError
	for _, validationErr := range validationErrors {
		if validationErr.Field == "profile.addresses[0].street" && validationErr.Tag == "required" {
			streetError = &validationErr
			break
		}
	}

	// Verify the custom message was applied correctly
	assert.NotNil(t, streetError, "Expected street validation error")
	if streetError != nil {
		assert.Equal(t, "Street is required", streetError.Message)
	}
}

func TestLeafNameFallbackLookup(t *testing.T) {
	// Create a validator
	validator := NewCustomValidator()

	// Set custom error message only for the leaf field name
	validator.translator.SetFieldError("street", "required", "Street address is required")

	// Create a test struct with nested fields
	type Address struct {
		Street string `json:"street" validate:"required"`
		City   string `json:"city" validate:"required"`
	}

	type User struct {
		HomeAddress Address `json:"homeAddress" validate:"required"`
		WorkAddress Address `json:"workAddress" validate:"required"`
	}

	// Create an invalid user with missing streets in both addresses
	invalidUser := User{
		HomeAddress: Address{
			Street: "", // Should get custom message from leaf field
			City:   "New York",
		},
		WorkAddress: Address{
			Street: "", // Should get same custom message from leaf field
			City:   "Chicago",
		},
	}

	// Validate
	err := validator.Validate(&invalidUser)
	assert.NotNil(t, err, "Expected validation error")

	// Convert to ValidationErrors
	validationErrors, ok := err.(ValidationErrors)
	assert.True(t, ok, "Expected ValidationErrors type")

	// Check both street validation errors
	homeStreetFound := false
	workStreetFound := false

	for _, validationErr := range validationErrors {
		if validationErr.Field == "homeAddress.street" && validationErr.Tag == "required" {
			homeStreetFound = true
			assert.Equal(t, "Street address is required", validationErr.Message)
		}
		if validationErr.Field == "workAddress.street" && validationErr.Tag == "required" {
			workStreetFound = true
			assert.Equal(t, "Street address is required", validationErr.Message)
		}
	}

	assert.True(t, homeStreetFound, "Expected home address street validation error")
	assert.True(t, workStreetFound, "Expected work address street validation error")
}
