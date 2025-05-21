package api

import (
	"net/http"

	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/labstack/echo/v4"
)

type apiErrorResponse struct {
	Errors map[string]string `json:"errors,omitempty"`
	apiResponse
}

func (err apiErrorResponse) Error() string {
	return err.Message
}

func errorHandler(err error, c echo.Context) {
	var (
		code                      = http.StatusInternalServerError
		message                   = "internal server error"
		errors  map[string]string = nil
	)

	switch err := err.(type) {
	case validator.ValidationErrors:
		code = http.StatusBadRequest
		message = "invalid request body"
		errors = make(map[string]string)
		for _, ve := range err {
			errors[ve.Path] = ve.Message
		}
	case validator.ValidationError:
		code = http.StatusBadRequest
		message = "invalid request body"
		errors = make(map[string]string)
		errors[err.Path] = err.Message
	case *echo.HTTPError:
		code = err.Code
		switch m := err.Message.(type) {
		case string:
			message = m
		case map[string]interface{}:
			if errMsg, ok := m["message"].(string); ok {
				message = errMsg
			}
		case error:
			message = m.Error()
		}
	case apiErrorResponse:
		code = err.Code
		message = err.Message
		errors = err.Errors
	}

	if !c.Response().Committed {
		err := c.JSON(code, apiErrorResponse{
			Errors: errors,
			apiResponse: apiResponse{
				Code:    code,
				Message: message,
			},
		})
		if err != nil {
			c.Logger().Errorf("Failed to respond after error: %v", err)
		}
	}
}
