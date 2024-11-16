package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"konbini/types"
)

// ErrorHandler handles all errors globally.
// Assign to echo.HTTPErrorHandler
func ErrorHandler(err error, c echo.Context) {
	var response *types.ErrorResponse

	if he, ok := err.(*echo.HTTPError); ok {
		if resp, ok := he.Message.(*types.ErrorResponse); ok {
			response = resp
		} else {
			response = &types.ErrorResponse{
				Status:  he.Code,
				Message: he.Message.(string),
			}
		}
	} else {
		// default error interface handling
		response = &types.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	// send error is no response sent yet
	if !c.Response().Committed {
		c.JSON(response.Status, response)
	}
}
