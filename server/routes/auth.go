package routes

import (
	commonApi "konbini/common/api"
	"konbini/server/handlers"
	"konbini/server/middlewares"
	"reflect"
)

func setupAuthRoutes(routeConfig *RouteConfig) {
	routeConfig.Echo.POST(
		commonApi.UriLogin,
		handlers.Login(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(commonApi.LoginRequest{})),
	)
	routeConfig.Echo.POST(
		commonApi.UriRegister,
		handlers.Register(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.RegisterRequest{})),
	)

	routeConfig.Echo.POST(
		commonApi.UriCheckToken,
		handlers.CheckAuthToken(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(commonApi.CheckAuthTokenRequest{})),
	)

	routeConfig.Echo.POST(
		commonApi.UriTOTPSetup,
		handlers.SetupTOTP(routeConfig.DBConnector),
		middlewares.ProtectAll(routeConfig.DBConnector),
	)
	routeConfig.Echo.POST(
		commonApi.UriTOTPLock,
		handlers.SetupTOTPLock(routeConfig.DBConnector),
		middlewares.ProtectAll(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(commonApi.SetupTOTPLockRequest{})),
	)
	routeConfig.Echo.DELETE(
		commonApi.UriTOTPDelete,
		handlers.RemoveTOTP(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.RemoveTOTPRequest{})),
	)

	routeConfig.Echo.GET(commonApi.UriVerifyEmail, handlers.VerifyEmail(routeConfig.DBConnector))
	routeConfig.Echo.POST(
		commonApi.UriResendVerificationEmail,
		handlers.ResendVerificationEmail(routeConfig.DBConnector),
		middlewares.ProtectAll(routeConfig.DBConnector),
	)
}
