package routes

import (
	"konbini/handlers"
)

func setupAuthRoutes(routeConfig *RouteConfig) {
	routeConfig.Echo.POST("/auth/register", handlers.HandleRegister())
}
