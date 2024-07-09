package router

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo/v4"
)

// writeJSON is a helper function that writes json to the client.
// The built-in c.JSON from echo does not set the content-length when writing the json response
// so this method sets it.
//
// IMPORTANT: This method should only be used for small json responses.
func writeJSON(status int, c echo.Context, i interface{}) error {
	payload, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(status)
	n, err := c.Response().Write(payload)
	c.Response().Header().Set(echo.HeaderContentLength, strconv.Itoa(n))
	return err
}
