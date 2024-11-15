package main

import (
	"konbini/common"
	"konbini/store"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
}
