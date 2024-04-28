package router

import (
	"bytes"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/juancwu/konbini/service"
	"github.com/juancwu/konbini/templates"
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
	Email        string `json:"email" validate:"required"`
	FirstName    string `json:"first_name" validate:"required,alpha"`
	LastName     string `json:"last_name" validate:"required,alpha"`
	PemPublicKey string `json:"pem_public_key" validate:"required"`
}

type VerifyEmailData struct {
	FirstName string
	LastName  string
	URL       string
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
	log.Info("POST /auth/register")
	reqBody := new(RegisterReqBody)

	// bind the incoming request data
	if err := c.Bind(reqBody); err != nil {
		log.Errorf("Error binding request body: %v\n", err)
		return c.String(http.StatusBadRequest, "Invalid payload")
	}

	if err := c.Validate(reqBody); err != nil {
		return err
	}

	user, err := service.GetUserByEmail(reqBody.Email)
	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Errorf("Error registering user: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	if user != nil {
		log.Error("Error registering user: an account with the given email already exists.", "email", reqBody.Email)
		return c.String(http.StatusBadRequest, "An account with the given email already exists.")
	}

	// register user
	userId, err := service.RegisterUser(reqBody.FirstName, reqBody.LastName, reqBody.Email, reqBody.PemPublicKey)
	if err != nil {
		log.Errorf("Error registering user: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	// create entry for an email verification
	err = service.CreateEmailVerification(userId)
	if err != nil {
		log.Errorf("Error creating email verificaiton: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error creating email verification.")
	}

	// get verify email template
	var tpl bytes.Buffer
	err = templates.Render(&tpl, "verify-email.html", VerifyEmailData{FirstName: reqBody.FirstName, LastName: reqBody.LastName, URL: "http://localhost:3000/verify-email"})
	if err != nil {
		log.Errorf("Failed to get verify email template: %v - handleRegister\n", err)
		return c.String(http.StatusInternalServerError, "Error sending verification email.")
	}

	log.Info("Request verify email", "func", "handleRegister")
	emailId, err := service.SendEmail("noreply@juancwu.dev", reqBody.Email, "[Konbini] Verify Your Email", tpl.String())
	log.Info("Verify email sent", "id", emailId, "func", "handleRegister")

	return c.String(http.StatusCreated, "Account registered.")
}
