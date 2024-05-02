package router

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/env"
	"github.com/juancwu/konbini/server/service"
	"github.com/juancwu/konbini/server/templates"
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
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"first_name" validate:"required,alpha"`
	LastName  string `json:"last_name" validate:"required,alpha"`
}

type VerifyEmailData struct {
	FirstName string
	LastName  string
	URL       string
}

func SetupAuthRoutes(e *echo.Echo) {
	e.POST("/auth", handleAuth)
	e.POST("/auth/register", handleRegister)
	e.GET("/auth/verify-email/:refId", handleVerifyEmail)
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
	userId, err := service.RegisterUser(reqBody.FirstName, reqBody.LastName, reqBody.Email)
	if err != nil {
		log.Errorf("Error registering user: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	// create entry for an email verification
	refId, err := service.CreateEmailVerification(userId)
	if err != nil {
		log.Errorf("Error creating email verificaiton: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error creating email verification.")
	}

	// get verify email template
	var tpl bytes.Buffer
	err = templates.Render(&tpl, "verify-email.html", VerifyEmailData{FirstName: reqBody.FirstName, LastName: reqBody.LastName, URL: fmt.Sprintf("%s/auth/verify-email/%s", env.Values().SERVER_URL, refId)})
	if err != nil {
		log.Errorf("Failed to get verify email template: %v - handleRegister\n", err)
		return c.String(http.StatusInternalServerError, "Error sending verification email.")
	}

	log.Info("Request verify email", "func", "handleRegister")
	emailId, err := service.SendEmail("noreply@juancwu.dev", reqBody.Email, "[Konbini] Verify Your Email", tpl.String())
	log.Info("Verify email sent", "id", emailId, "func", "handleRegister")

	return c.String(http.StatusCreated, "Account registered.")
}

func handleVerifyEmail(c echo.Context) error {
	log.Info("GET /auth/verify-email/:refId")
	refId := c.Param("refId")
	if refId == "" {
		log.Info("Invalid request to verify email when no ref id was found.")
		return c.String(http.StatusBadRequest, "Missing value.")
	}

	ev, err := service.GetEmailVerification(refId)
	if err != nil {
		log.Errorf("Error verifying email: %v\n", err)
		return c.String(http.StatusInternalServerError, "Could not verify email. Please try again later.")
	}

	// update the user entry that email has been verified
	log.Info("Updating user entry to set email_verified...")
	_, err = database.DB().Exec("UPDATE users SET email_verified = true WHERE id = $1;", ev.UserId)
	if err != nil {
		log.Errorf("Error updating user entry to set email_verified to true: %v\n", err)
		return c.String(http.StatusInternalServerError, "Could not verify email. Please try again later.")
	}

	// now we can update the email verification status because user entry has been updated
	log.Info("Updating email verification status...")
	_, err = database.DB().Exec("UPDATE email_verifications SET status = $1 WHERE id = $2;", service.EMAIL_STATUS_VERIFIED, ev.Id)
	if err != nil {
		// this error doesn't matter that much as long as the user entry has been updated
		log.Errorf("Error updating email verification entry to set status to verified: %v\n", err)
	}

	log.Info("Email verified")
	return c.String(http.StatusOK, "Email verified.")
}
