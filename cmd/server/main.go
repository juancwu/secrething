package main

import (
	"konbini/server/config"
	"konbini/server/db"
	"konbini/server/handlers"
	"konbini/server/routes"
	inner_validator "konbini/server/validator"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load server configuration")
	}

	dbUrl, dbAuthToken := cfg.GetDatabaseConfig()

	e := echo.New()

	validate := validator.New()
	cv := inner_validator.Validator{Validator: validate}
	e.Validator = &cv

	// set global error handler
	e.HTTPErrorHandler = handlers.ErrorHandler()

	// v1 routes
	apiV1 := e.Group("/api/v1")
	routeConfig := &routes.RouteConfig{
		Echo:         apiV1,
		ServerConfig: cfg,
		DBConnector:  db.NewConnector(dbUrl, dbAuthToken),
	}
	routes.SetupRoutesV1(routeConfig)

	err = e.Start(cfg.GetPort())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server.")
	}
}
