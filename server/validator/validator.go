package validator

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator is a wrapper around the validator.Validate instance
type CustomValidator struct {
	validator  *govalidator.Validate
	translator *ErrorTranslator
}

// ErrorTranslator handles translating validation errors to custom messages
type ErrorTranslator struct {
	fieldErrors    map[string]map[string]string // map[field]map[tag]message
	defaultErrors  map[string]string            // map[tag]message
	defaultMessage string
}

// ValidationContext stores context-specific validation settings
type ValidationContext struct {
	validator  *CustomValidator
	translator *ErrorTranslator
}

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
func (v ValidationErrors) AsMap() map[string]string {
	errMap := make(map[string]string)
	for _, err := range v {
		errMap[err.Field] = err.Message
	}
	return errMap
}

// FormatValidationErrors formats ValidationErrors into a standardized map structure
// This is used by error handlers to format validation errors consistently
func FormatValidationErrors(valErrors ValidationErrors) map[string]interface{} {
	fieldErrors := make(map[string]interface{})

	for _, validationErr := range valErrors {
		fieldErrors[validationErr.Field] = validationErr.Message
	}

	return fieldErrors
}

// NewErrorTranslator creates a new ErrorTranslator
func NewErrorTranslator() *ErrorTranslator {
	return &ErrorTranslator{
		fieldErrors:    make(map[string]map[string]string),
		defaultErrors:  make(map[string]string),
		defaultMessage: "Invalid value",
	}
}

// SetFieldError sets a custom error message for a specific field and validation tag
func (t *ErrorTranslator) SetFieldError(field, tag, message string) {
	if _, ok := t.fieldErrors[field]; !ok {
		t.fieldErrors[field] = make(map[string]string)
	}
	t.fieldErrors[field][tag] = message
}

// SetDefaultError sets a default error message for a validation tag
func (t *ErrorTranslator) SetDefaultError(tag, message string) {
	t.defaultErrors[tag] = message
}

// SetDefaultMessage sets the default error message for all validations
func (t *ErrorTranslator) SetDefaultMessage(message string) {
	t.defaultMessage = message
}

// Translate translates a validation error to a custom message
func (t *ErrorTranslator) Translate(field string, tag string) string {
	// Check if there's a custom message for this field and tag
	if fieldMessages, ok := t.fieldErrors[field]; ok {
		if message, ok := fieldMessages[tag]; ok {
			return message
		}
	}

	// Check if there's a default message for this tag
	if message, ok := t.defaultErrors[tag]; ok {
		return message
	}

	// Return the default message
	return t.defaultMessage
}

// setDefaultMessages sets commonly used validation error messages
func setDefaultMessages(translator *ErrorTranslator) {
	defaults := map[string]string{
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
	}

	for tag, message := range defaults {
		translator.SetDefaultError(tag, message)
	}
}

// NewCustomValidator creates a new CustomValidator instance
func NewCustomValidator() *CustomValidator {
	v := govalidator.New()
	// Register function to get field name from json tag
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	translator := NewErrorTranslator()
	setDefaultMessages(translator)

	return &CustomValidator{
		validator:  v,
		translator: translator,
	}
}

// Validate validates the provided struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		validationErrors := ValidationErrors{}

		for _, err := range err.(govalidator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()

			// Get the JSON field name
			structField, _ := reflect.TypeOf(i).Elem().FieldByName(err.StructField())
			jsonName := strings.SplitN(structField.Tag.Get("json"), ",", 2)[0]
			if jsonName != "" {
				field = jsonName
			}

			// Translate the error
			message := cv.translator.Translate(field, tag)

			// Create a validation error
			validationError := ValidationError{
				Field:   field,
				Message: message,
				Tag:     tag,
				Value:   err.Value(),
			}

			validationErrors = append(validationErrors, validationError)
		}

		return validationErrors
	}

	return nil
}

// Translator gets the CustomValidator's ErrorTranslator instance
func (cv *CustomValidator) Translator() *ErrorTranslator {
	return cv.translator
}

// Clone creates a copy of the ErrorTranslator
func (t *ErrorTranslator) Clone() *ErrorTranslator {
	clone := NewErrorTranslator()
	clone.defaultMessage = t.defaultMessage

	// Copy default errors
	for tag, msg := range t.defaultErrors {
		clone.defaultErrors[tag] = msg
	}

	// Copy field errors
	for field, tagMsgs := range t.fieldErrors {
		clone.fieldErrors[field] = make(map[string]string)
		for tag, msg := range tagMsgs {
			clone.fieldErrors[field][tag] = msg
		}
	}

	return clone
}

// Clone creates a copy of the CustomValidator with a new translator
func (cv *CustomValidator) Clone() *CustomValidator {
	return &CustomValidator{
		validator:  cv.validator,
		translator: cv.translator.Clone(),
	}
}

// NewValidationContext creates a new validation context with a cloned validator
func NewValidationContext(baseValidator *CustomValidator) *ValidationContext {
	cloned := baseValidator.Clone()
	return &ValidationContext{
		validator:  cloned,
		translator: cloned.translator,
	}
}

// SetFieldError sets a custom error message for a specific field and validation tag
func (vc *ValidationContext) SetFieldError(field, tag, message string) *ValidationContext {
	vc.translator.SetFieldError(field, tag, message)
	return vc
}

// SetDefaultError sets a default error message for a validation tag
func (vc *ValidationContext) SetDefaultError(tag, message string) *ValidationContext {
	vc.translator.SetDefaultError(tag, message)
	return vc
}

// SetDefaultMessage sets the default error message for all validations
func (vc *ValidationContext) SetDefaultMessage(message string) *ValidationContext {
	vc.translator.SetDefaultMessage(message)
	return vc
}

// Validate validates the provided struct using this context's validator
func (vc *ValidationContext) Validate(i interface{}) error {
	return vc.validator.Validate(i)
}

// BindAndValidate binds and validates a request body to a struct
func BindAndValidate(c echo.Context, i interface{}) error {
	// Bind the request body to the struct
	if err := c.Bind(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate the struct
	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}

// BindAndValidateWithContext binds and validates a request body to a struct using a custom validation context
func BindAndValidateWithContext(c echo.Context, i interface{}, vc *ValidationContext) error {
	// Bind the request body to the struct
	if err := c.Bind(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate the struct using the context's validator
	if err := vc.Validate(i); err != nil {
		return err
	}

	return nil
}
