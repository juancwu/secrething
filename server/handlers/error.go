package handlers

import (
	"fmt"
	commonApi "konbini/common/api"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type APIError struct {
	Code          int      `json:"code"`
	PublicMessage string   `json:"message"`
	Errors        []string `json:"errors,omitempty"`
	RequestId     string   `json:"request_id"`

	PrivateMessage string `json:"-"`
	InternalError  error  `json:"-"`
}

func (e APIError) Error() string {
	return fmt.Sprintf(
		"API Error - Code: %d - Public Msg: %s - Private Msg: %s - Internal Error: %s - Request ID: %s",
		e.Code,
		e.PublicMessage,
		e.PrivateMessage,
		e.InternalError.Error(),
		e.RequestId,
	)
}

func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var apiError APIError
		switch err.(type) {
		case *echo.HTTPError:
			he := err.(*echo.HTTPError)
			apiError.Code = he.Code
			apiError.PublicMessage = he.Error()
			apiError.InternalError = he.Internal
		case APIError:
			apiError = err.(APIError)
		default:
			apiError.InternalError = err
			apiError.Code = http.StatusInternalServerError
			apiError.PublicMessage = http.StatusText(apiError.Code)
			apiError.PrivateMessage = err.Error()
		}

		requestId := c.Request().Header.Get(echo.HeaderXRequestID)
		apiError.RequestId = requestId
		path := c.Request().URL.Path
		method := c.Request().Method
		ip := c.RealIP()

		log.
			Error().
			Err(apiError.InternalError).
			Str("request_id", apiError.RequestId).
			Int("code", apiError.Code).
			Str("path", path).
			Str("method", method).
			Str("ip", ip).
			Str("public_message", apiError.PublicMessage).
			Msg(apiError.PrivateMessage)

		if !c.Response().Committed {
			err := c.JSON(apiError.Code, commonApi.ErrorResponse{
				Code:      apiError.Code,
				Message:   apiError.PublicMessage,
				Errors:    apiError.Errors,
				RequestId: apiError.RequestId,
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to response from error handler.")
			}
		}
	}
}
