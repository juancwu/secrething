package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/juancwu/konbini/server/db"
	"github.com/juancwu/konbini/server/infrastructure/middleware"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load .env")
	}

	// Establish database connection
	tursoConnector := db.NewTursoConnector(os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_AUTH_TOKEN"))
	_, err := tursoConnector.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to establish connection with database")
	}

	e := echo.New()
	e.HTTPErrorHandler = middleware.ErrorHandlerMiddleware()

	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server.")
	}
}
