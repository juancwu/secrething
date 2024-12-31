package handlers

import (
	"errors"
	"konbini/server/middlewares"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterRequest represents the request body for register route.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=32"`
	NickName string `json:"nickname" validate:"required,min=3,max=32"`
}

// HandleRegister is a handler function that registers a user for Konbini.
func HandleRegister() echo.HandlerFunc {
	return func(c echo.Context) error {
		body, ok := c.Get(middlewares.JSON_BODY_KEY).(*RegisterRequest)
		if !ok {
			return errors.New("Failed to get JSON body from context.")
		}

		return c.JSON(http.StatusOK, body)
	}
}
