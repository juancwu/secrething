package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// apiError represents a way to define a custom apiError that ErrHandler can identify and use.
type apiError struct {
	// Code is the http response status code.
	Code int `json:"-"`
	// Err is the original error. This gets log by ErrHandler.
	// This field can also be a validator.ValidationErrs.
	Err error `json:"-"`
	// Msg is the internal message that gets log.
	Msg string `json:"-"`

	// Errs is a string that gets send back to the client after ValidationErrs is processed. This field will be in the json response if not empty.
	Errs []string `json:"errors,omitempty"`
	// PublicMsg is the message that gets send back to the client. This field will be in the json response if not empty.
	PublicMsg string `json:"message,omitempty"`
	// RequestId is the id of the request the error came from, helps with traceback. This field will be in the json response if not empty.
	RequestId string `json:"request_id,omitempty"`
}

// Error satisfies the error interface so that apiError can be used as an error type.
func (e apiError) Error() string {
	return e.Msg
}

// ErrHandler is a custom error handler that will log the error and corresponding message.
// Use an echo.HTTPError if there is a need to return a status code other than 500.
// Normal errors will be handled using a 500 and generic internal server error message..
func ErrHandler(err error, c echo.Context) {
	switch err.(type) {
	case *echo.HTTPError:
		he := err.(*echo.HTTPError)
		log.Error().Err(he).Msg("Echo HTTPError")
		writeJSON(he.Code, c, map[string]string{"message": http.StatusText(http.StatusInternalServerError)})
	case apiError:
		e := err.(apiError)
		if e.Code == 0 {
			e.Code = http.StatusInternalServerError
		}

		if ve, ok := e.Err.(validator.ValidationErrors); ok {
			if e.PublicMsg == "" {
				e.PublicMsg = "Invalid request body. Please fix the issues."
			}
			e.Errs = make([]string, len(ve))
			re, err := regexp.Compile(`\[(\d+)\]`)
			if err != nil {
				log.Error().Err(err).Str(echo.HeaderXRequestID, e.RequestId).Msg("Failed to compile regex.")
				re = nil
			}
			for i, err := range ve {
				field := fmt.Sprintf("%s.%s", err.StructNamespace(), err.Tag())
				var sliceIndeces []string
				if re != nil {
					matches := re.FindAllStringSubmatch(field, -1)
					for _, m := range matches {
						sliceIndeces = append(sliceIndeces, m[1])
					}
					field = re.ReplaceAllString(field, "")
				}
				log.Info().Msg(field)
				msg, exists := reqBodyValidationMsgs[field]
				if !exists {
					msg = fmt.Sprintf("Validation failed on the '%s' failed.", err.Tag())
				}
				if len(sliceIndeces) > 0 {
					msg = fmt.Sprintf("%s (index path: %s)", msg, strings.Join(sliceIndeces, "->"))
				}
				e.Errs[i] = msg
			}
		}

		if e.PublicMsg == "" {
			e.PublicMsg = http.StatusText(e.Code)
		}

		log.Error().
			Err(e.Err).
			Str(echo.HeaderXRequestID, e.RequestId).
			Int("status_code", e.Code).
			Msg(e.Msg)

		writeJSON(e.Code, c, e)
	default:
		log.Error().Msg("Standard error encountered. Somewhere in route code is returning standard error.")
		c.NoContent(http.StatusInternalServerError)
	}
}
