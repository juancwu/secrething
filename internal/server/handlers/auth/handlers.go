package auth

import (
	"context"
	"net/http"
	"time"

	handlerErrors "github.com/juancwu/secrething/internal/server/handlers/errors"
	"github.com/juancwu/secrething/internal/server/middleware"
	authService "github.com/juancwu/secrething/internal/server/services/auth"
	"github.com/labstack/echo/v4"
)

// createUser handles requests to create new users.
func createUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Minute)
	defer cancel()

	requestID := ""

	var body createUserRequest
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := middleware.Validate(c, &body); err != nil {
		return err
	}

	user, err := authService.CreateUser(ctx, body.Email, body.Password, body.Name)
	if err != nil {
		if serviceErr, ok := err.(authService.AuthServiceError); ok && serviceErr.Is(authService.UserAlreadyExistsErr) {
			return handlerErrors.NewBadRequest(serviceErr.Error(), "", UserAlreadyExistsCode, requestID, err)
		}
	}

	return c.JSON(http.StatusOK, createUserResponse{UserID: user.UserID})
}

// signinUser handles requests to sign-in users.
func signinUser(c echo.Context) error {
	return nil
}
