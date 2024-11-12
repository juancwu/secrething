package router

import (
	"net/http"

	"github.com/juancwu/konbini/config"
	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// SetupHealthcheckRoutes set ups all the routes to do healthcheck.
func SetupHealthcheckRoutes(e RouterGroup) {
	e.GET("/health", handleGetHealth)
}

// handleGetHealth handles basic healthcheck requests.
func handleGetHealth(c echo.Context) error {
	report := healthReport{
		Database: true,
		Version:  config.Version,
	}
	// ping the db
	if err := store.Ping(); err != nil {
		log.Error().Err(err).Str(echo.HeaderXRequestID, c.Request().Header.Get(echo.HeaderXRequestID)).Msg("Database failed healthcheck.")
		report.Database = false
	}
	return writeJSON(http.StatusOK, c, report)
}
