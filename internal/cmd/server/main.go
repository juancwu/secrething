package main

import (
	"github.com/juancwu/konbini/internal/server/config"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := config.Load(".env"); err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = !config.IsDevelopment()

	if err := e.Start(config.Server().Address); err != nil {
		panic(err)
	}
}
