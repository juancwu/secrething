package middleware

import (
	"errors"
	"konbini/types"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var validate = validator.New()

const (
	// constant key to retrieve the validated interface from echo.Context
	REQUEST_MODEL_CTX_KEY = "mw_bind_validate_request_model"
)

// BindAndValidate binds and validates the request body
func BindAndValidate(structType reflect.Type) echo.MiddlewareFunc {
	if structType == nil {
		// Crash, invalid state
		log.Panic().Err(errors.New("Nil struct type passed to BindAndValidate")).Stack().Msg("CRASH: invalid state reached.")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// create a new copy for the request model using the given structType
			model := reflect.New(structType)
			i := model.Interface()
			if err := c.Bind(i); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, &types.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "Invalid request body",
				})
			}

			if err := validate.StructCtx(c.Request().Context(), i); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, types.NewValidationError(structType, err))
			}

			c.Set(REQUEST_MODEL_CTX_KEY, i)

			return next(c)
		}
	}
}
