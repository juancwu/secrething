package router

import "github.com/labstack/echo/v4"

// RouterGroup is a simple interface to allow passing a echo.Echo or echo.Group instance.
// This makes it easier to extend the router with or without nested routes.
type RouteGroup interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// HealthReport is a collection of stats about the backend.
type HealthReport struct {
	// DB can be "health" or "unhealthy"
	DB string `json:"db"`
}
