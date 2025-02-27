package routes

import "github.com/juancwu/konbini/server/handlers"

// setupHealthRoutes sets all the health related routes
func setupHealthRoutes(routeConfig *RouteConfig) {
	routeConfig.Echo.GET("/health-check", handlers.HealthCheck(routeConfig.ServerConfig, routeConfig.DBConnector))
}
