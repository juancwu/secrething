package router

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// err_msg_logger_key is the key for the field in echo.Context to set/get the custom error message.
// Keep in mind that this is a custom defined constant and it has nothing to do with echo.
const (
	err_msg_logger_key = "err_msg"
	public_err_msg_key = "public_err_msg"
)

// errorResponseBody represents the response body when an api error occurs
type errorResponseBody struct {
	Message   string   `json:"message"`
	RequestId string   `json:"request_id,omitempty"`
	Errors    []string `json:"errors,omitempty"`
}

// ErrHandler is a custom error handler that will log the error and corresponding message.
// Use an echo.HTTPError if there is a need to return a status code other than 500.
// Normal errors will be handled using a 500 and generic internal server error message..
func ErrHandler(err error, c echo.Context) {
	publicErrMsg, ok := c.Get(public_err_msg_key).(string)
	if !ok || publicErrMsg == "" {
		publicErrMsg = "internal server error"
	}

	requestId, ok := c.Get(echo.HeaderXRequestID).(string)
	if !ok {
		requestId = ""
	}

	body := errorResponseBody{
		Message:   publicErrMsg,
		RequestId: requestId,
	}

	code := http.StatusInternalServerError

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	} else if errs, ok := err.(validator.ValidationErrors); ok {
		// format the validation error message
		body.Message = "Invalid Request Body"
		errMsgs := make([]string, len(errs))
		for i, err := range errs {
			field := fmt.Sprintf("%s.%s", err.StructNamespace(), err.Tag())
			msg, exists := reqBodyValidationMsgs[field]
			if !exists {
				msg = fmt.Sprintf("Validation failed on the '%s' failed.", err.Tag())
			}
			errMsgs[i] = msg
		}
		body.Errors = errMsgs
		code = http.StatusBadRequest
	}

	errMsg, ok := c.Get(err_msg_logger_key).(string)
	if !ok || errMsg == "" {
		errMsg = "internal server error"
	}

	log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Int("status_code", code).Msg(errMsg)

	c.JSON(
		code,
		body,
	)
}
