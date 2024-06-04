package router

import (
	"time"

	"github.com/labstack/echo/v4"
)

type RouteGroup interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type ApiError struct {
	StatusCode int    `json:"status_code"`
	Msg        string `json:"msg"`
	RequestId  string `json:"request_id"`
}

type ApiReqBodyError struct {
	StatusCode int    `json:"status_code"`
	Errors     any    `json:"errors"`
	RequestId  string `json:"request_id"`
}

type RequestBodyValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type RebrandBentoRequestBody struct {
	NewName string `json:"new_name" validate:"required,min=3,max=255,ascii"`
}

// GetChallangeResp represents the response body that will be sent back
// when a client requests a new challenge.
type GetChallangeResp struct {
	State     string    `json:"state"`
	Challange string    `json:"challenge"`
	ExpiresAt time.Time `json:"expires_at"`
}
