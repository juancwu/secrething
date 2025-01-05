package routes

import (
	"konbini/server/handlers"
	"konbini/server/middlewares"
	"reflect"
)

func setupAuthRoutes(routeConfig *RouteConfig) {
	routeConfig.Echo.POST(
		"/auth/login",
		handlers.Login(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.LoginRequest{})),
	)
	routeConfig.Echo.POST(
		"/auth/register",
		handlers.Register(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.RegisterRequest{})),
	)

	// routeConfig.Echo.POST(
	// 	"/auth/magic/request",
	// 	handlers.HandleMagicLinkRequest(routeConfig.DBConnector),
	// 	middlewares.ValidateJson(reflect.TypeOf(handlers.MagicLinkRequestRequest{})),
	// )
	// routeConfig.Echo.GET(
	// 	"/auth/magic/verify",
	// 	handlers.HandleMagicLinkVerify(routeConfig.DBConnector),
	// )
	// routeConfig.Echo.POST(
	// 	"/auth/magic/status",
	// 	handlers.HandleMagicLinkStatus(),
	// 	middlewares.ValidateJson(reflect.TypeOf(handlers.MagicLinkStatusRequest{})),
	// )

	routeConfig.Echo.GET("/auth/email/verify", handlers.VerifyEmail(routeConfig.DBConnector))
}
