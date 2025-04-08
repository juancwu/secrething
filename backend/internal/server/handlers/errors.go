package handlers

import (
	"fmt"

	"github.com/juancwu/secrething/internal/server/api"
	"github.com/labstack/echo/v4"
)

func ErrorHandler() echo.HTTPErrorHandler {
	return handleError
}

func handleError(err error, c echo.Context) {
	switch err := err.(type) {
	case api.AppError:
		fmt.Println(err.Err)
		c.JSON(err.ResponseCode, err)
	default:
		c.JSON(500, map[string]interface{}{"error": err.Error()})
	}
}
