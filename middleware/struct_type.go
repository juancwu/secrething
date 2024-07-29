package middleware

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

const (
	STRUCT_TYPE_KEY string = "structype"
)

// StructType will register a struct type in the echo.Context for the incoming request.
// Optionally, provide a key for the struct type. Pass empty string to use default key "structype".
func StructType(structType reflect.Type) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(STRUCT_TYPE_KEY, structType)
			return next(c)
		}
	}
}
