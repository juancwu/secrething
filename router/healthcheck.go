package router

import (
	"net/http"
	"os"

	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
)

// SetupHealthcheckRoutes sets the routes for healthcheck
func SetupHealthcheckRoutes(e RouteGroup) {
	e.GET("/health", handleGetHealth)
}

// GetHealth gets a simple overview of the current health of the backend.
func handleGetHealth(c echo.Context) error {
	report := HealthReport{
		DB:      "healthy",
		Version: os.Getenv("VERSION"),
	}
	err := store.Ping()
	if err != nil {
		report.DB = "unhealthy"
	}
	return c.JSON(http.StatusOK, report)
}
