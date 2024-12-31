package routes

import "konbini/handlers"

// setupHealthRoutes sets all the health related routes
func setupHealthRoutes(routeConfig *RouteConfig) {
	routeConfig.Echo.GET("/health-check", handlers.HandleHealthCheck(routeConfig.ServerConfig, routeConfig.DatabaseConn))
}
