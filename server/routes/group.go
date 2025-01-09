package routes

import (
	"konbini/server/handlers"
	"konbini/server/middlewares"
	"reflect"
)

func setupGroupRoutes(routeConfig *RouteConfig) {
	e := routeConfig.Echo

	e.POST(
		"/group/new",
		handlers.NewGroup(routeConfig.DBConnector),
		// only allow request that comes with a full token
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.NewGroupRequest{})),
	)
}
