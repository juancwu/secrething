package router

import "github.com/labstack/echo/v4"

type RouterGroup interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// basicRespBody represents a basic JSON response body that only has a "message" field.
type basicRespBody struct {
	Msg       string `json:"message"`
	RequestId string `json:"request_id,omitempty"`
}

// healthReport represents the report after doing a healthcheck on the server.
type healthReport struct {
	Database bool `json:"database"`
}
