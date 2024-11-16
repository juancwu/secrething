package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"konbini/common"
	"konbini/handler"
	"konbini/store"
)

func main() {
	zerolog.TimeFieldFormat = common.FRIENDLY_TIME_FORMAT
	if os.Getenv("APP_ENV") == common.DEVELOPMENT_ENV {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: common.FRIENDLY_TIME_FORMAT})
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("Failed to load .env")
		}
	}

	db, err := store.NewConn()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load .env")
	}
	defer db.Close()

	e := echo.New()

	e.HTTPErrorHandler = handler.ErrorHandler
}
