package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	e := echo.New()

	err := e.Start(os.Getenv("PORT"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server.")
	}
}
