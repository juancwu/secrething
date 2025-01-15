package routes

import (
	"konbini/server/handlers"
	"konbini/server/middlewares"
	"reflect"
)

func setupBentoRoutes(routeConfig *RouteConfig) {
	e := routeConfig.Echo

	e.GET(
		"/bento",
		handlers.GetBento(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
	)
	e.GET(
		"/bentos",
		handlers.ListBentos(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
	)
	e.GET(
		"/bento/metadata",
		handlers.GetBentoMetadata(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
	)
	e.POST(
		"/bento/new",
		handlers.NewBento(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.NewBentoRequest{})),
	)
	e.POST(
		"/bento/ingredients",
		handlers.AddIngredientsToBento(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.AddIngredientsToBentoRequest{})),
	)
	e.DELETE(
		"/bento/ingredients",
		handlers.RemoveIngredientsFromBento(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.RemoveIngredientsFromBentoRequest{})),
	)
}
