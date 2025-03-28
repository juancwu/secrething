package middleware

import (
	"errors"

	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/labstack/echo/v4"
)

type Registry interface {
	GetMessages() validator.ValidationMessages
}

func SetValidator(v *validator.Validator, r Registry) echo.MiddlewareFunc {
	customV := v.UseMessages(r.GetMessages())
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("validator", customV)
			return next(c)
		}
	}
}

func Validate(c echo.Context, i interface{}) error {
	v, ok := c.Get("validator").(*validator.Validator)
	if !ok {
		return errors.New("Validation attempt without registering validator")
	}
	return v.Validate(i)
}
