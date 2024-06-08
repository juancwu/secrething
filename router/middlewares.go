// This file contains all the middlewares used within the router.
package router

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// useValidateRequestBody is a middleware that given a struct type, it will validate it
// using the validator that was setup when creating a new echo.Echo
func useValidateRequestBody(i interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			body := reflect.New(reflect.TypeOf(i)).Interface()
			if err := c.Bind(body); err != nil {
				logger.Error(
					"Failed to bind request body",
					zap.String("method", c.Request().Method),
					zap.String("path", c.Request().URL.Path),
					zap.Error(err),
				)
				return c.JSON(
					http.StatusInternalServerError,
					apiResponse{
						StatusCode: http.StatusInternalServerError,
						Message: map[string]string{
							"errors": http.StatusText(http.StatusInternalServerError),
						},
					},
				)
			}
			if err := c.Validate(body); err != nil {
				var ve validator.ValidationErrors
				if errors.As(err, &ve) {
					return c.JSON(
						http.StatusBadRequest,
						apiResponse{
							StatusCode: http.StatusBadRequest,
							Message: map[string]string{
								"errors": ve.Error(),
							},
						},
					)
				}

				logger.Error(
					"Failed to request body",
					zap.String("method", c.Request().Method),
					zap.String("path", c.Request().URL.Path),
					zap.Error(err),
				)
				return c.JSON(
					http.StatusInternalServerError,
					apiResponse{
						StatusCode: http.StatusInternalServerError,
						Message: map[string]string{
							"errors": http.StatusText(http.StatusInternalServerError),
						},
					},
				)
			}

			c.Set("body", body)
			return next(c)
		}
	}
}
