package main

import (
	"os"

	"github.com/juancwu/konbini/config"
	"github.com/juancwu/konbini/router"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// faster time field format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// load env, checks for required env
	if os.Getenv("APP_ENV") == config.DEV_ENV {
		err := config.LoadEnv()
		if err != nil {
			log.Panic().Err(err).Msg("Failed to load env.")
		}
	}

	// setup echo
	e := echo.New()
	e.HTTPErrorHandler = router.ErrHandler

	// start echo
	err := e.Start(":" + os.Getenv("PORT"))
	if err != nil {
		log.Panic().Err(err).Msg("Failed to start echo server.")
	}
}
