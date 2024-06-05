package router

import (
	"fmt"
	"net/http"
	"strings"

	// package modules
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"

	// local modules
	"github.com/juancwu/konbini/server/service"
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
			logger, _ := zap.NewProduction()
			defer logger.Sync()
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
					logger.Error("Validate request error", zap.Error(err))
					return writeApiError(c, http.StatusInternalServerError, "internal server error")
				}
				if value == "" && !validator.Required {
					return next(c)
				}
				if value == "" && validator.Required {
					return writeApiError(c, http.StatusBadRequest, fmt.Sprintf("Missing required %s parameter \"%s\"", validator.From, validator.Field))
				}
				if validator.Validate != nil {
					err = validator.Validate(value)
					if err != nil {
						return writeApiError(c, http.StatusBadRequest, err.Error())
					}
				}
			}
			return next(c)
		}
	}
}

// RequestID is a middleware that attaches a randomly generated request id to
// the request header with echo.HeaderXRequestID as key.
func RequestID(l int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId, err := gonanoid.New(32)
			if err != nil {
				logger, _ := zap.NewProduction()
				defer logger.Sync()
				logger.Error("Error generating request id", zap.Error(err))
			}
			c.Request().Header.Add(echo.HeaderXRequestID, requestId)

			return next(c)
		}
	}
}

// Logger is a middleware that logs basic information of incoming requests.
// It will log out the method, path and request id.
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger, _ := zap.NewProduction()
			defer logger.Sync()

			logger.Info(
				"New incoming request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("request_id", c.Request().Header.Get(echo.HeaderXRequestID)),
			)

			return next(c)
		}
	}
}

func JwtAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or malformed Bearer token")
		}

		token, err := service.VerifyToken(parts[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if claims, ok := token.Claims.(*service.JwtCustomClaims); ok && token.Valid {
			c.Set("token", token)
			c.Set("claims", claims)
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}
}
