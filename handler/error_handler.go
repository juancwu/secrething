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
			var msg string
			if m, ok := he.Message.(string); ok {
				msg = m
			} else {
				msg = http.StatusText(he.Code)
			}
			response = &types.ErrorResponse{
				Status:  he.Code,
				Message: msg,
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
