package main

import (
	"github.com/juancwu/go-valkit/v2/validations"
	"github.com/juancwu/secrething/internal/server/config"
	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/handlers"
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
	e.HTTPErrorHandler = handlers.ErrorHandler()

	authHandler := handlers.NewAuthHandler()
	authHandler.ConfigureRoutes(e, v)

	if err := e.Start(config.Server().Address); err != nil {
		panic(err)
	}
}
