package routes

import (
	"konbini/server/handlers"
	"konbini/server/middlewares"
	"reflect"
)

func setupBentoRoutes(routeConfig *RouteConfig) {
	e := routeConfig.Echo

	e.POST(
		"/bento/new",
		handlers.NewBento(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.NewBentoRequest{})),
	)
}
