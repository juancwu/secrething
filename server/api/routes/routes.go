package routes

import (
	"github.com/juancwu/konbini/server/api/handlers/auth"
	"github.com/juancwu/konbini/server/config"
	"github.com/juancwu/konbini/server/db"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(e *echo.Echo, cfg *config.Config, db *db.TursoConnector) {
	// Auth routes
	authHandler := auth.NewAuthHandler(cfg, db)

	// Define auth routes
	authGroup := e.Group("/api/v1/auth")
	authGroup.POST("/register", authHandler.Register)

	// Add more route groups here as they are implemented
}
