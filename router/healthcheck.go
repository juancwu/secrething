package router

import (
	"net/http"

	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
)

// SetupHealthcheckRoutes sets the routes for healthcheck
func SetupHealthcheckRoutes(e RouteGroup) {
	e.GET("/health", getHealth)
}

// GetHealth gets a simple overview of the current health of the backend.
func getHealth(c echo.Context) error {
	report := HealthReport{
		DB: "healthy",
	}
	err := store.Ping()
	if err != nil {
		report.DB = "unhealhty"
	}
	return c.JSON(http.StatusOK, report)
}
