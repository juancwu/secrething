package main

import (
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/juancwu/konbini/server/api/routes"
	"github.com/juancwu/konbini/server/config"
	"github.com/juancwu/konbini/server/db"
	"github.com/juancwu/konbini/server/middleware"
	"github.com/juancwu/konbini/server/observability"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Establish database connection
	tursoConnector := db.NewTursoConnector(cfg.DatabaseURL, cfg.DatabaseAuthToken)
	_, err = tursoConnector.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to establish connection with database")
	}

	err = observability.InitSentry(observability.SentryConfig{
		DSN:              cfg.SentryDSN,
		Environment:      string(cfg.Environment),
		Debug:            cfg.Debug,
		SampleRate:       1.0,
		TracesSampleRate: 0.2,
		MaxBreadcrumbs:   100,
		EnableTracing:    true,
		ServerName:       cfg.ServerName,
	})

	e := echo.New()

	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.CSRF())
	e.Use(middleware.LoggerMiddleware())

	// Sentry setup
	e.Use(sentryecho.New(sentryecho.Options{}))
	e.Use(observability.SentryHubMiddleware())

	// Global HTTP error handler
	e.HTTPErrorHandler = middleware.ErrorHandlerMiddleware()

	// Register routes
	routes.RegisterRoutes(e, cfg, tursoConnector)

	if err := e.Start(cfg.GetAddress()); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server.")
	}
}
