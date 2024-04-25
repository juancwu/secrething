package router

import (
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/juancwu/konbini/service"
	"github.com/labstack/echo/v4"
)

/*
   All the routes in here are prefixed with /auth
*/

type AuthReqBody struct {
	Email     string `json:"email" validate:"required"`
	Challenge string `json:"challenge" validate:"required"`
}

type RegisterReqBody struct {
	Email        string  `json:"email" validate:"required"`
	FirstName    *string `json:"first_name" validate:"alpha"`
	LastName     *string `json:"last_name" validate:"alpha"`
	PemPublicKey string  `json:"pem_public_key" validate:"required"`
}

func SetupAuthRoutes(e *echo.Echo) {
	e.POST("/auth", handleAuth)
	e.POST("/auth/register", handleRegister)
}

func handleAuth(c echo.Context) error {
	auth := new(AuthReqBody)

	// bind the incoming request data
	if err := c.Bind(auth); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(auth); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, auth)
}

func handleRegister(c echo.Context) error {
	reqBody := new(RegisterReqBody)

	// bind the incoming request data
	if err := c.Bind(reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(reqBody); err != nil {
		return err
	}

	user, err := service.GetUserByEmail(reqBody.Email)
	if err != nil {
		log.Errorf("Error registering user: %v\n", err)
		c.Response().WriteHeader(http.StatusInternalServerError)
		return errors.New("Error registering user.")
	}

	if user != nil {
		log.Error("Error registering user: an account with the given email already exists.", reqBody.Email)
		c.Response().WriteHeader(http.StatusBadRequest)
		return errors.New("An account with the given email already exists.")
	}

	// register user
	err = service.RegisterUser(reqBody.FirstName, reqBody.LastName, reqBody.Email, reqBody.PemPublicKey)
	if err != nil {
		log.Errorf("Error registering user: %v\n", err)
		c.Response().WriteHeader(http.StatusInternalServerError)
		return errors.New("Error registering user.")
	}

	c.Response().WriteHeader(http.StatusCreated)

	return nil
}
