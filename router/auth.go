package router

import "github.com/labstack/echo/v4"

/*
   All the routes in here are prefixed with /auth
*/

func SetupAuthRouter(e *echo.Echo) {
	// route that registers a user
	e.POST("/auth/register", handleRegister)
}

func handleRegister(c echo.Context) error {
	return nil
}
