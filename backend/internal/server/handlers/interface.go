package handlers

import (
	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	ConfigureRoutes(e *echo.Echo, v *validator.Validator)
}
