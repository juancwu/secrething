package validator

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
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
func (v ValidationErrors) AsMap() map[string]interface{} {
	return FormatValidationErrors(v)
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
	// First, check if there's a custom message for the full field path and tag
	if fieldMessages, ok := t.fieldErrors[field]; ok {
		if message, ok := fieldMessages[tag]; ok {
			return message
		}
	}

	// If no full path match, check if there's a leaf field name match
	// Only apply this if the field contains a dot (meaning it's a nested field)
	if strings.Contains(field, ".") {
		leafField := field
		if idx := strings.LastIndex(field, "."); idx >= 0 {
			leafField = field[idx+1:]
		}

		if fieldMessages, ok := t.fieldErrors[leafField]; ok {
			if message, ok := fieldMessages[tag]; ok {
				return message
			}
		}
	}

	// Check if there's a default message for this tag
	if message, ok := t.defaultErrors[tag]; ok {
		return message
	}

	// Return the default message
	return t.defaultMessage
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

	// Register custom validation tag for passwords
	v.RegisterValidation("password", validatePassword)

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
			tag := err.Tag()

			// Get the JSON field name path
			jsonFieldPath := getJSONFieldPath(i, err)

			// Translate the error using the full path - our enhanced translator
			// will handle both full path and leaf name lookups
			message := cv.translator.Translate(jsonFieldPath, tag)

			// Create a validation error
			validationError := ValidationError{
				Field:   jsonFieldPath,
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

// getJSONFieldPath returns the JSON field path for a validation error
// For nested fields, it returns a dot-separated path like "profile.firstName"
// For array/slice fields, it returns indexed paths like "items[0].name"
func getJSONFieldPath(obj interface{}, fieldError govalidator.FieldError) string {
	// Build the namespace path based on JSON field names rather than struct field names
	namespace := fieldError.Namespace()
	parts := strings.Split(namespace, ".")

	// The first part is the type name, so we skip it
	parts = parts[1:]

	// Build a new path with JSON names
	var jsonParts []string
	currentObj := obj
	currentType := reflect.TypeOf(currentObj).Elem()
	currentValue := reflect.ValueOf(currentObj).Elem()

	for i, part := range parts {
		// Check if this part refers to an array/slice index
		indexMatch := regexp.MustCompile(`^(\w+)\[(\d+)\]$`).FindStringSubmatch(part)

		if len(indexMatch) == 3 {
			// This is an array/slice index reference like 'Items[0]'
			fieldName := indexMatch[1]
			indexStr := indexMatch[2]
			index, _ := strconv.Atoi(indexStr)

			// Find the struct field for the array/slice
			field, found := currentType.FieldByName(fieldName)
			if !found {
				// If we can't find the field, just use the original part
				jsonParts = append(jsonParts, part)
				continue
			}

			// Get the JSON tag name for the array/slice field
			jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if jsonName == "" || jsonName == "-" {
				jsonName = fieldName
			}

			// Add the field name and index to the path
			jsonParts = append(jsonParts, fmt.Sprintf("%s[%d]", jsonName, index))

			// Update currentType/currentValue for the next iteration
			fieldValue := currentValue.FieldByName(fieldName)
			if !fieldValue.IsValid() || index >= fieldValue.Len() {
				// If the field value is invalid or index is out of bounds, we can't continue
				// Just append the remaining parts as is
				for j := i + 1; j < len(parts); j++ {
					jsonParts = append(jsonParts, parts[j])
				}
				break
			}

			// Get the element at the specified index
			elemValue := fieldValue.Index(index)

			// Update currentType and currentValue based on the element type
			if elemValue.Kind() == reflect.Struct {
				currentType = elemValue.Type()
				currentValue = elemValue
			} else if elemValue.Kind() == reflect.Ptr && elemValue.Elem().Kind() == reflect.Struct {
				currentType = elemValue.Elem().Type()
				currentValue = elemValue.Elem()
			} else {
				// If the element is not a struct, we can't go deeper
				// Just append the remaining parts as is
				for j := i + 1; j < len(parts); j++ {
					jsonParts = append(jsonParts, parts[j])
				}
				break
			}
		} else {
			// Regular struct field
			field, found := currentType.FieldByName(part)
			if !found {
				// If we can't find the field, just use the original part
				jsonParts = append(jsonParts, part)
				continue
			}

			// Get the JSON tag name
			jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if jsonName == "" || jsonName == "-" {
				// If there's no JSON tag or it's "-", use the original field name
				jsonParts = append(jsonParts, part)
			} else {
				jsonParts = append(jsonParts, jsonName)
			}

			// Update currentType and currentValue for the next iteration
			fieldValue := currentValue.FieldByName(part)

			if field.Type.Kind() == reflect.Struct {
				currentType = field.Type
				if fieldValue.IsValid() {
					currentValue = fieldValue
				}
			} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
				currentType = field.Type.Elem()
				if fieldValue.IsValid() && !fieldValue.IsNil() {
					currentValue = fieldValue.Elem()
				}
			} else {
				// If the field is not a struct or pointer to struct, we can't go deeper
				// Just append the remaining parts as is
				for j := i + 1; j < len(parts); j++ {
					jsonParts = append(jsonParts, parts[j])
				}
				break
			}
		}
	}

	return strings.Join(jsonParts, ".")
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
