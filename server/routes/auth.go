package routes

import (
	"konbini/server/handlers"
	"konbini/server/middlewares"
	"reflect"
)

func setupAuthRoutes(routeConfig *RouteConfig) {
	routeConfig.Echo.POST(
		"/auth/register",
		handlers.HandleRegister(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.RegisterRequest{})),
	)
	routeConfig.Echo.GET("/auth/email/verify", handlers.VerifyEmail(routeConfig.DBConnector))
}
