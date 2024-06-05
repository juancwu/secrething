package router

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/juancwu/konbini/server/utils"
)

const (
	AUTH_ROUTER_PREFIX = "/auth"
)

// getInvalidTagMsg gets a generic error message for the given tag.
func getInvalidTagMsg(tag string) string {
	switch tag {
	case "required":
		return "This field is required."
	}
	return "Invalid field"
}

// writeApiError is a utility function that writes a json response with ApiError.
func writeApiError(c echo.Context, statusCode int, msg string) error {
	return c.JSON(statusCode, ApiError{StatusCode: statusCode, Msg: msg, RequestId: c.Request().Header.Get(echo.HeaderXRequestID)})
}

// writeApiReqBodyError is a utility function that writes a json response with ApiReqBodyError
func writeApiReqBodyError(c echo.Context, statusCode int, ve validator.ValidationErrors) error {
	out := make([]RequestBodyValidationError, len(ve))
	for i, fe := range ve {
		out[i] = RequestBodyValidationError{
			Field:  fe.Field(),
			Reason: getInvalidTagMsg(fe.Tag()),
		}
	}
	return c.JSON(http.StatusBadRequest, ApiReqBodyError{
		StatusCode: http.StatusBadRequest,
		Errors:     out,
		RequestId:  c.Request().Header.Get(echo.HeaderXRequestID),
	})
}

func writeNoBody(c echo.Context, statusCode int) error {
	c.Response().Writer.WriteHeader(statusCode)
	c.Response().Writer.Header().Add(echo.HeaderXRequestID, c.Request().Header.Get(echo.HeaderXRequestID))
	return nil
}

func validateUUID(uuid string) error {
	if utils.IsValidUUIDV4(uuid) {
		return nil
	}
	return fmt.Errorf("The given id is not a proper UUID v4: %s", uuid)
}
