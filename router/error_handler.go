package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// err_msg_logger_key is the key for the field in echo.Context to set/get the custom error message.
// Keep in mind that this is a custom defined constant and it has nothing to do with echo.
const err_msg_logger_key = "err_msg"

// ErrHandler is a custom error handler that will log the error and corresponding message.
// Use an echo.HTTPError if there is a need to return a status code other than 500.
// Normal errors will be handled using a 500 and generic internal server error message..
func ErrHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	requestId, ok := c.Get(echo.HeaderXRequestID).(string)
	if !ok {
		requestId = "nil"
	}

	msg, ok := c.Get(err_msg_logger_key).(string)
	if !ok || msg == "" {
		msg = "internal server error"
	}

	log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Int("status_code", code).Msg(msg)

	c.JSON(
		code,
		map[string]string{
			"message":    "internal server error",
			"request_id": requestId,
		},
	)
}
