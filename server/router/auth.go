package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/env"
	usermodel "github.com/juancwu/konbini/server/models/user"
	"github.com/juancwu/konbini/server/service"
	"github.com/juancwu/konbini/server/templates"
	"github.com/juancwu/konbini/server/utils"
)

/*
   All the routes in here are prefixed with /auth
*/

type AuthReqBody struct {
	Email     string `json:"email" validate:"required"`
	Challenge string `json:"challenge" validate:"required"`
}

type RegisterReqBody struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required,alpha"`
	LastName  string `json:"last_name" validate:"required,alpha"`
	Password  string `json:"password" validate:"required,min=12"`
}

type VerifyEmailData struct {
	FirstName string
	LastName  string
	URL       string
}

type ResetPasswordEmailData struct {
	FirstName string
	LastName  string
	ResetId   string
}

func SetupAuthRoutes(e *echo.Echo) {
	e.POST("/auth", handleAuth)
	e.POST("/auth/register", handleRegister)
	e.GET("/auth/verify-email/:refId", handleVerifyEmail)
	e.POST("/auth/reset/password", handleStartResetPassword)
	e.PATCH("/auth/reset/password", handleFinishResetPassword)
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
		utils.Logger().Errorf("Error binding request body: %v\n", err)
		return c.String(http.StatusBadRequest, "Invalid payload")
	}

	if err := c.Validate(reqBody); err != nil {
		return err
	}

	user, err := service.GetUserWithEmail(reqBody.Email)
	if err != nil && err.Error() != "sql: no rows in result set" {
		utils.Logger().Errorf("Error registering user: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	if user != nil {
		utils.Logger().Error("Error registering user: an account with the given email already exists.", "email", reqBody.Email)
		return c.String(http.StatusBadRequest, "An account with the given email already exists.")
	}

	tx, err := database.DB().Begin()
	if err != nil {
		utils.Logger().Errorf("Error beginning transaction to register user: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	// register user
	userId, err := service.RegisterUser(reqBody.FirstName, reqBody.LastName, reqBody.Email, reqBody.Password, tx)
	if err != nil {
		utils.Logger().Errorf("Error registering user: %v\n", err)
		database.Rollback(tx, c.Request().URL.Path)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	// create entry for an email verification
	refId, err := service.CreateEmailVerification(userId, tx)
	if err != nil {
		utils.Logger().Errorf("Error creating email verificaiton: %v\n", err)
		database.Rollback(tx, c.Request().URL.Path)
		return c.String(http.StatusInternalServerError, "Error registering user")
	}

	// want to commit before sending the email so that we have a record of the email verification
	err = tx.Commit()
	if err != nil {
		utils.Logger().Errorf("Error committing transaction changes (%s): %v\n", c.Request().URL.Path, err)
		database.Rollback(tx, c.Request().URL.Path)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	// get verify email template
	var tpl bytes.Buffer
	err = templates.Render(&tpl, "verify-email.html", VerifyEmailData{FirstName: reqBody.FirstName, LastName: reqBody.LastName, URL: fmt.Sprintf("%s/auth/verify-email/%s", env.Values().SERVER_URL, refId)})
	if err != nil {
		utils.Logger().Errorf("Failed to get verify email template: %v - handleRegister\n", err)
		return c.String(http.StatusInternalServerError, "Error sending verification email.")
	}

	utils.Logger().Info("Request verify email", "func", "handleRegister")
	emailId, err := service.SendEmail(env.Values().NOREPLY_EMAIL, reqBody.Email, "[Konbini] Verify Your Email", tpl.String())
	utils.Logger().Info("Verify email sent", "id", emailId, "func", "handleRegister")

	_, err = database.DB().Exec("UPDATE email_verifications SET email_sent_at = $1, resend_email_id = $2, status = $3 WHERE verification_id = $4;", time.Now().In(time.UTC), emailId, service.EMAIL_STATUS_SENT, refId)
	if err != nil {
		utils.Logger().Errorf("Failed to update email verification with sent time and resend email id: %v\n", err)
	}

	return c.String(http.StatusCreated, "Account created! Please verify your email to unlock all Konbini services.")
}

func handleVerifyEmail(c echo.Context) error {
	utils.Logger().Info("GET /auth/verify-email/:refId")
	routeErrorMessage := "Could not verify email. Please try again later."
	refId := c.Param("refId")
	if refId == "" {
		utils.Logger().Info("Invalid request to verify email when no ref id was found.")
		return c.String(http.StatusBadRequest, "Missing value.")
	}

	ev, err := service.GetEmailVerification(refId)
	if err != nil {
		utils.Logger().Errorf("Error verifying email: %v\n", err)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	// update the user entry that email has been verified
	utils.Logger().Info("Updating user entry to set email_verified...")
	tx, err := database.DB().Begin()
	if err != nil {
		utils.Logger().Errorf("Error begining transaction to verify email: %v\n", err)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}
	_, err = tx.Exec("UPDATE users SET email_verified = true WHERE id = $1;", ev.UserId)
	if err != nil {
		utils.Logger().Errorf("Error updating user entry to set email_verified to true: %v\n", err)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	// now we can update the email verification status because user entry has been updated
	utils.Logger().Info("Updating email verification status...")
	_, err = tx.Exec("UPDATE email_verifications SET status = $1, verified_at = $2 WHERE id = $3;", service.EMAIL_STATUS_VERIFIED, time.Now().In(time.UTC), ev.Id)
	if err != nil {
		// this error doesn't matter that much as long as the user entry has been updated
		utils.Logger().Errorf("Error updating email verification entry to set status to verified: %v\n", err)
		database.Rollback(tx, c.Request().URL.Path)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	err = tx.Commit()
	if err != nil {
		utils.Logger().Errorf("Error committing verify email changes: %v\n", err)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	utils.Logger().Info("Email verified")
	return c.String(http.StatusOK, "Email verified.")
}

func handleStartResetPassword(c echo.Context) error {
	routeErrorMessage := "Could not start reset password process. Please try again later."

	// limit read to 1024 bytes or characters, more than enough for a single email
	utils.Logger().Info("Getting email from request body")
	bodyReader := io.LimitReader(c.Request().Body, 1024)
	reqBody, err := io.ReadAll(bodyReader)
	if err != nil {
		utils.Logger().Errorf("Failed to read request body: %v\n", err)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}
	email := string(reqBody)

	user, err := usermodel.GetByEmail(email)

	// generate random reset id
	utils.Logger().Info("Generating reset id...")
	resetId, err := gonanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 12)
	if err != nil {
		utils.Logger().Errorf("Failed to generate reset id: %v\n", err)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}
	utils.Logger().Infof("Reset id: %s\n", resetId)

	// expires in 1 hour
	expTime := time.Now().In(time.UTC).Add(time.Hour * 1)
	utils.Logger().Infof("Expire time reset password: %s\n", expTime.String())

	// transaction
	utils.Logger().Info("Starting transaction to store reset password record")
	tx, err := database.DB().Begin()
	if err != nil {
		utils.Logger().Errorf("Failed to start transaction to store reset password record: %v\n", err)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	// check if there is another record of the user already, if there is, delete that one
	utils.Logger().Info("Checking if reset record exists for given user...")
	var recordId int64 = 0
	err = database.DB().QueryRow("SELECT id FROM users_passwords_resets WHERE user_id = $1;", user.Id).Scan(&recordId)
	if err != nil && err.Error() != "sql: no rows in result set" {
		utils.Logger().Error("Failed to check reset record: %v\n", err)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	if recordId != 0 {
		utils.Logger().Info("Updating existing password record...")
		_, err = tx.Exec("UPDATE users_passwords_resets SET reset_id = $1, expires_at = $2 WHERE id = $3;", resetId, expTime, recordId)
		if err != nil {
			database.Rollback(tx, "reset password")
			return c.String(http.StatusInternalServerError, routeErrorMessage)
		}
	} else {
		utils.Logger().Info("Saving new reset password record")
		_, err = tx.Exec("INSERT INTO users_passwords_resets (user_id, reset_id, expires_at) VALUES ($1, $2, $3);", user.Id, resetId, expTime)
		if err != nil {
			utils.Logger().Errorf("Failed to insert password reset link id into db: %v\n", err)
			database.Rollback(tx, "reset password")
			return c.String(http.StatusInternalServerError, routeErrorMessage)
		}

	}

	err = tx.Commit()
	if err != nil {
		utils.Logger().Errorf("Failed to commit changes from transcation: %v\n", err)
		database.Rollback(tx, "reset password")
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}
	utils.Logger().Info("Reset id saved", "reset_id", resetId)

	// send email with reset id
	utils.Logger().Info("Generating email template...")
	var tpl bytes.Buffer
	err = templates.Render(&tpl, "reset-password.html", ResetPasswordEmailData{FirstName: user.FirstName, LastName: user.LastName, ResetId: resetId})
	utils.Logger().Info("Sending email with reset id for password reset...")
	_, err = service.SendEmail(env.Values().NOREPLY_EMAIL, email, "[Konbini] Reset Your Password", tpl.String())
	if err != nil {
		utils.Logger().Errorf("Failed to send password reset email: %v\n", err)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Code: %s", resetId))
}

func handleFinishResetPassword(c echo.Context) error {
	return c.String(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}
