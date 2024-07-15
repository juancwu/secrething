package router

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	c.Response().Header().Set(echo.HeaderContentLength, strconv.Itoa(len(payload)))
	c.Response().WriteHeader(status)
	_, err = c.Response().Write(payload)
	return err
}

// readRequestBody reads a request body and validates it.
func readRequestBody(c echo.Context, body interface{}) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Binding request body.")
	if err := c.Bind(body); err != nil {
		return err
	}
	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Validating request body.")
	if err := c.Validate(body); err != nil {
		return err
	}
	return nil
}
