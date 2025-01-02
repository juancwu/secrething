package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

const (
	JSON_BODY_KEY string = "request_validated_json"
)

var (
	ErrFailedToGetJsonBody error = errors.New("Failed to get json body from context")
)

type ValidateJsonConfig struct {
	MaxBodySize int64
}

func ValidateJson(structType reflect.Type) echo.MiddlewareFunc {
	return ValidateJsonWithConfig(structType, ValidateJsonConfig{
		MaxBodySize: 10 << 20, // 10MB
	})
}

func ValidateJsonWithConfig(structType reflect.Type, cfg ValidateJsonConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contentType := c.Request().Header.Get(echo.HeaderContentType)
			if contentType != echo.MIMEApplicationJSON {
				return echo.NewHTTPError(http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			}

			// limit the request body size
			c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, cfg.MaxBodySize)

			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				if err.Error() == "http: request body too large" {
					he := echo.NewHTTPError(http.StatusRequestEntityTooLarge, fmt.Sprintf("The request body is too large. The maximum size is %d bytes.", cfg.MaxBodySize))
					he.SetInternal(err)
					return he
				}
				return err
			}

			// restore the body for later use
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

			if len(bytes.TrimSpace(body)) == 0 {
				return echo.NewHTTPError(http.StatusBadRequest, "Request body must not be empty.")
			}

			if !json.Valid(body) {
				return echo.NewHTTPError(http.StatusBadRequest, "Request body is not a valid JSON.")
			}

			reqStruct := reflect.New(structType)
			i := reqStruct.Interface()

			// read in the request body and make sure it is within size limit
			err = json.Unmarshal(body, i)
			if err != nil {
				return err
			}

			if err := c.Validate(i); err != nil {
				// this will let the global error handler handle
				// the ValidationError and get error string for
				// the each invalid field.
				return err
			}

			// allow the remaining handlers in the chain gain access to
			// the request body.
			c.Set(JSON_BODY_KEY, i)

			return next(c)
		}
	}
}

func GetJsonBody[T interface{}](c echo.Context) (*T, error) {
	body, ok := c.Get(JSON_BODY_KEY).(*T)
	if !ok {
		return nil, ErrFailedToGetJsonBody
	}
	return body, nil
}
