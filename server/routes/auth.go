package routes

import (
	"konbini/server/handlers"
	"konbini/server/middlewares"
	"reflect"
)

func setupAuthRoutes(routeConfig *RouteConfig) {
	routeConfig.Echo.POST(
		"/auth/register",
		handlers.HandleRegister(routeConfig.Queries),
		middlewares.ValidateJson(reflect.TypeOf(handlers.RegisterRequest{})),
	)
}
