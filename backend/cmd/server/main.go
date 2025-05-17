package main

import (
	"github.com/joho/godotenv"
	"github.com/juancwu/secrething/internal/config"
	"github.com/juancwu/secrething/internal/db"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load .env file first
	godotenv.Load()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	e := echo.New()

	// Pass config to DB connect
	conn, err := db.Connect(cfg)
	if err != nil {
		e.Logger.Fatal(err)
	}

	conn.Ping()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	if err := e.Start(cfg.Server.Address); err != nil {
		e.Logger.Fatal(err)
	}
}
