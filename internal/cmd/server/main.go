package main

import (
	"github.com/juancwu/go-valkit/v2/validations"
	"github.com/juancwu/secrething/internal/server/config"
	"github.com/juancwu/secrething/internal/server/db"
	authHandler "github.com/juancwu/secrething/internal/server/handlers/auth"
	"github.com/juancwu/secrething/internal/server/handlers/errors"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := config.Load(".env"); err != nil {
		panic(err)
	}

	if _, err := db.Connect(); err != nil {
		panic(err)
	}

	v := config.DefaultValidator()
	validations.AddPasswordValidation(v, validations.DefaultPasswordOptions())

	e := echo.New()
	e.HideBanner = !config.IsDevelopment()
	e.HTTPErrorHandler = errors.ErrorHandler()

	apiGroup := e.Group("/api")
	authGroup := apiGroup.Group("/auth")

	authHandler.Configure(authGroup, v)

	if err := e.Start(config.Server().Address); err != nil {
		panic(err)
	}
}
