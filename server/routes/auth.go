package routes

import (
	commonApi "github.com/juancwu/konbini/common/api"
	"github.com/juancwu/konbini/server/handlers"
	"github.com/juancwu/konbini/server/middlewares"
	"reflect"

	"github.com/labstack/echo/v4"
)

func setupAuthRoutes(routeConfig *RouteConfig) {
	// Login route with rate limiting for TOTP verification
	loginRoute := routeConfig.Echo.Group(commonApi.UriLogin)
	loginRoute.Use(middlewares.ValidateJson(reflect.TypeOf(commonApi.LoginRequest{})))
	loginRoute.POST("", handlers.Login(routeConfig.DBConnector))
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
	// Setup TOTP Lock route with rate limiting
	totpLockRoute := routeConfig.Echo.Group(commonApi.UriTOTPLock)
	totpLockRoute.Use(
		middlewares.ProtectAll(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(commonApi.SetupTOTPLockRequest{})),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				user, err := middlewares.GetUser(c)
				if err != nil {
					return err
				}
				// Apply rate limiting middleware with user ID
				return middlewares.LimitTOTPAttempts(user.ID)(next)(c)
			}
		},
	)
	totpLockRoute.POST("", handlers.SetupTOTPLock(routeConfig.DBConnector))
	// TOTP Delete route with rate limiting
	totpDeleteRoute := routeConfig.Echo.Group(commonApi.UriTOTPDelete)
	totpDeleteRoute.Use(
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.RemoveTOTPRequest{})),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				user, err := middlewares.GetUser(c)
				if err != nil {
					return err
				}
				// Apply rate limiting middleware with user ID
				return middlewares.LimitTOTPAttempts(user.ID)(next)(c)
			}
		},
	)
	totpDeleteRoute.DELETE("", handlers.RemoveTOTP(routeConfig.DBConnector))

	routeConfig.Echo.GET(commonApi.UriVerifyEmail, handlers.VerifyEmail(routeConfig.DBConnector))
	routeConfig.Echo.POST(
		commonApi.UriResendVerificationEmail,
		handlers.ResendVerificationEmail(routeConfig.DBConnector),
		middlewares.ProtectAll(routeConfig.DBConnector),
	)
}
