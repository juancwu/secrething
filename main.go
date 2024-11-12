package main

import (
	"konbini/common"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = common.FRIENDLY_TIME_FORMAT
	if os.Getenv("APP_ENV") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: common.FRIENDLY_TIME_FORMAT})
	}
}
