package routes

import (
	"database/sql"
	"konbini/server/config"
	"konbini/server/db"

	"github.com/labstack/echo/v4"
)

// EchoInstance interface represents the Echo struct implementation which allows setup functions
// to take in echo.Echo or echo.Group instead of just one of the two options.
type EchoInstance interface {
	Use(middleware ...echo.MiddlewareFunc)
	CONNECT(path string, h echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	Group(prefix string, middleware ...echo.MiddlewareFunc) *echo.Group
}

// RouteConfig represents the configuration that is needed to setup any route.
type RouteConfig struct {
	Echo         EchoInstance
	ServerConfig *config.Config
	DatabaseConn *sql.DB
	Queries      *db.Queries
}
