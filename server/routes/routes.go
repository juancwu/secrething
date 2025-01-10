package routes

import (
	"konbini/server/middlewares"

	"github.com/rs/zerolog/log"
)

// SetupRoutes is a global setup function for all the routes.
func SetupRoutesV1(cfg *RouteConfig) {
	cfg.Echo.Use(middlewares.LoggerWithConfig(
		middlewares.LoggerConfig{
			Logger:  log.Logger,
			Exclude: []string{},
		},
	))
	setupAuthRoutes(cfg)
	setupHealthRoutes(cfg)
	setupGroupRoutes(cfg)
	setupBentoRoutes(cfg)
}
