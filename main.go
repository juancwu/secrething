package main

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/juancwu/konbini/config"
	"github.com/juancwu/konbini/router"
	"github.com/juancwu/konbini/store"
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

	// connect to database
	err := store.Connect(os.Getenv("DB_URL"))
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to database.")
	}

	// setup echo
	e := echo.New()
	e.HTTPErrorHandler = router.ErrHandler

	validate := validator.New()
	validate.RegisterValidation("password", validatePassword)
	cv := customValidator{validator: validate}
	e.Validator = &cv
	// start echo
	err = e.Start(":" + os.Getenv("PORT"))
	if err != nil {
		log.Panic().Err(err).Msg("Failed to start echo server.")
	}
}
