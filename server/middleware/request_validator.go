package middleware

import (
	"fmt"
	"net/http"

	"github.com/juancwu/konbini/server/utils"
	"github.com/labstack/echo/v4"
)

const (
	VALIDATE_HEADER = "header"
	VALIDATE_QUERY  = "query"
	VALIDATE_PARAM  = "param"
)

// ValidatorFunc is called with the given a value to allow flexible
// way to validate the value of a query parameter. The function returns an error that
// will be used in the response back to the client.
type ValidatorFunc func(string) error

// ValidatorOptions represents the configuration of query parameter validator.
type ValidatorOptions struct {
	// Field is the name of the query/request param or header to validate.
	Field string
	// Required is a modifier that indicates whether the field being validated is required.
	// Empty strings will result in an error.
	Required bool
	// Validator is the ValidatorFunc that will be called to validate the field.
	Validate ValidatorFunc
	// From decides from where the field will be gotten from to perform validation.
	From string
}

// ValidateRequest is a middleware that accepts a series of ValidatorOptions
// that defines the parameters that need validation before a route handler runs.
func ValidateRequest(validators ...ValidatorOptions) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			for _, validator := range validators {
				var value string
				switch validator.From {
				case VALIDATE_QUERY:
					value = c.QueryParam(validator.Field)
				case VALIDATE_PARAM:
					value = c.Param(validator.Field)
				case VALIDATE_HEADER:
					value = c.Request().Header.Get(validator.Field)
				default:
					err = fmt.Errorf("Invalid ValidatorOptions.Field value. Expected one of %s, %s, or %s but received %s\n", VALIDATE_QUERY, VALIDATE_PARAM, VALIDATE_HEADER, validator.Field)
					utils.Logger().Error(err)
					return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}
				if value == "" && !validator.Required {
					return next(c)
				}
				if value == "" && validator.Required {
					return c.String(http.StatusBadRequest, fmt.Sprintf("Missing required %s parameter \"%s\"", validator.From, validator.Field))
				}
				if validator.Validate != nil {
					err = validator.Validate(value)
					if err != nil {
						return c.String(http.StatusBadRequest, err.Error())
					}
				}
			}
			return next(c)
		}
	}
}
