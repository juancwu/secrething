package validator

import (
	"net/http"
	"reflect"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Validator is a wrapper around the validator.Validate instance
type Validator struct {
	validator *govalidator.Validate
	resolver  *ErrorResolver
}

// NewValidator creates a new CustomValidator instance
func NewValidator() *Validator {
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

	return &Validator{
		validator: v,
		resolver:  translator,
	}
}

// Validate validates the provided struct
func (cv *Validator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		validationErrors := ValidationErrors{}

		for _, err := range err.(govalidator.ValidationErrors) {
			tag := err.Tag()

			// Get the JSON field name path
			jsonFieldPath := getJSONFieldPath(i, err)

			// Translate the error using the full path - our enhanced translator
			// will handle both full path and leaf name lookups
			message := cv.resolver.Translate(jsonFieldPath, tag)

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

// Resolver gets the CustomValidator's ErrorTranslator instance
func (cv *Validator) Resolver() *ErrorResolver {
	return cv.resolver
}

// Clone creates a copy of the CustomValidator with a new translator
func (cv *Validator) Clone() *Validator {
	return &Validator{
		validator: cv.validator,
		resolver:  cv.resolver.Clone(),
	}
}

// SetFieldError sets a custom error message for a specific field and validation tag
func (v *Validator) SetFieldError(field, tag, message string) *Validator {
	v.resolver.SetFieldError(field, tag, message)
	return v
}

// SetDefaultError sets a default error message for a validation tag
func (v *Validator) SetDefaultError(tag, message string) *Validator {
	v.resolver.SetDefaultError(tag, message)
	return v
}

// SetDefaultMessage sets the default error message for all validations
func (v *Validator) SetDefaultMessage(message string) *Validator {
	v.resolver.SetDefaultMessage(message)
	return v
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
func BindAndValidateWithContext(c echo.Context, i interface{}, vc *Validator) error {
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
