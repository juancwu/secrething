package handlers

import (
	"context"
	"database/sql"
	"konbini/server/config"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// HealthReport represents a health check report response body.
type HealthReport struct {
	Version                  string `json:"version"`
	DatabaseConnectionStatus string `json:"database_connection_status"`
}

// HandleHealthCheck handles health check requests.
// It gets the current running version of the app.
// It gets the database connection status.
func HandleHealthCheck(cfg *config.Config, conn *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		report := HealthReport{
			Version: cfg.GetVersion(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		err := conn.PingContext(ctx)
		dbStatus := "Healthy"
		if err != nil {
			log.Error().Err(err).Msg("Failed to ping database during health check")
			dbStatus = "Error"
		}

		report.DatabaseConnectionStatus = dbStatus

		return c.JSON(http.StatusOK, report)
	}
}
