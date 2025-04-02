package errors

import "github.com/labstack/echo/v4"

func ErrorHandler() echo.HTTPErrorHandler {
	return handleError
}

func handleError(err error, c echo.Context) {}
