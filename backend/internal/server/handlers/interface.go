package handlers

import (
	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	ConfigureRoutes(g *echo.Group, v *validator.Validator)
}
