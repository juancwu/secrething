package router

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/env"
	usermodel "github.com/juancwu/konbini/server/models/user"
	"github.com/juancwu/konbini/server/service"
	"github.com/juancwu/konbini/server/templates"
)

/*
   All the routes in here are prefixed with /auth
*/

type AuthReqBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
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
	ResetId   string `json:"reset_id" validate:"required,len=12"`
	Password  string `json:"password" validate:"required,min=12"`
}

func SetupAccountRoutes(e RouteGroup) {
	e.POST("/account/login", handleLogin)
	e.POST("/account/signup", handleSignup)
	e.GET("/account/verify-email", handleVerifyEmail, ValidateRequest(
		ValidatorOptions{
			Field:    "code",
			From:     VALIDATE_QUERY,
			Required: true,
			Validate: func(s string) error {
				return nil
			},
		},
	))
	e.POST("/account/reset-password", handleStartResetPassword)
	e.PATCH("/account/password", handleFinishResetPassword)
}

func handleLogin(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	auth := new(AuthReqBody)

	// bind the incoming request data
	if err := c.Bind(auth); err != nil {
		logger.Error("Error binding request body", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Authentication service down. Please try again later.")
	}
	if err := c.Validate(auth); err != nil {
		logger.Error("Authentication request body validation failed", zap.Error(err))
		return c.String(http.StatusBadRequest, "Bad request")
	}

	// find user
	logger.Info("Getting user information for authentication")
	user, err := usermodel.GetByEmailWithPassword(auth.Email, auth.Password)
	if err != nil {
		logger.Error("Failed to get user for authentication", zap.Error(err))
		if err == sql.ErrNoRows {
			return c.String(http.StatusBadRequest, "Invalid credentials. Make sure you inputted the right email and password.")
		}
		return c.String(http.StatusInternalServerError, "Authentication service down. Please try again later.")
	}

	// generate
	logger.Info("Generating access and refresh tokens")
	accessToken, err := service.GenerateAccessToken(user.Id)
	if err != nil {
		logger.Error("Failed to generate access token", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Authentication service down. Please try again later.")
	}

	return c.String(http.StatusOK, accessToken)
}

func handleSignup(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	reqBody := new(RegisterReqBody)

	// bind the incoming request data
	if err := c.Bind(reqBody); err != nil {
		logger.Error("Error binding request body", zap.Error(err))
		return c.String(http.StatusBadRequest, "Invalid payload")
	}

	if err := c.Validate(reqBody); err != nil {
		return err
	}

	user, err := usermodel.GetByEmail(reqBody.Email)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("Error registering user", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	if user != nil {
		logger.Error("Error registering user: an account with the given email already exists.", zap.String("email", reqBody.Email))
		return c.String(http.StatusBadRequest, "An account with the given email already exists.")
	}

	tx, err := database.DB().Begin()
	if err != nil {
		logger.Error("Error beginning transaction to register user", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	// register user
	logger.Info("Registering user with email", zap.String("email", reqBody.Email))
	row := tx.QueryRow(
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, crypt($4, gen_salt($5))) RETURNING id;",
		reqBody.FirstName, reqBody.LastName, reqBody.Email, reqBody.Password, env.Values().PASS_ENCRYPT_ALGO)
	if row.Err() != nil {
		logger.Error("Error resgitering user", zap.Error(row.Err()))
		database.Rollback(tx, c.Request().URL.Path)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	var userId string
	err = row.Scan(&userId)
	if err != nil {
		logger.Error("Error getting returning user id after insert", zap.Error(row.Err()))
		database.Rollback(tx, c.Request().URL.Path)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	// create entry for an email verification
	refId, err := service.CreateEmailVerification(userId, tx)
	if err != nil {
		logger.Error("Error creating email verificaiton", zap.Error(row.Err()))
		database.Rollback(tx, c.Request().URL.Path)
		return c.String(http.StatusInternalServerError, "Error registering user")
	}

	// want to commit before sending the email so that we have a record of the email verification
	err = tx.Commit()
	if err != nil {
		logger.Error("Error committing transaction changes", zap.Error(err))
		database.Rollback(tx, c.Request().URL.Path)
		return c.String(http.StatusInternalServerError, "Error registering user.")
	}

	// get verify email template
	if env.Values().APP_ENV == env.PRODUCTION {
		var tpl bytes.Buffer
		err = templates.Render(&tpl, "verify-email.html", VerifyEmailData{FirstName: reqBody.FirstName, LastName: reqBody.LastName, URL: fmt.Sprintf("%s/auth/verify-email/%s", env.Values().SERVER_URL, refId)})
		if err != nil {
			logger.Error("Failed to get verify email template", zap.Error(err))
			return c.String(http.StatusInternalServerError, "Error sending verification email.")
		}

		logger.Info("Sending verification email")
		emailId, err := service.SendEmail(env.Values().NOREPLY_EMAIL, reqBody.Email, "[Konbini] Verify Your Email", tpl.String())
		logger.Info("Verify email sent", zap.String("email_id", emailId))
		_, err = database.DB().Exec("UPDATE email_verifications SET email_sent_at = $1, resend_email_id = $2, status = $3 WHERE verification_id = $4;", time.Now().In(time.UTC), emailId, service.EMAIL_STATUS_SENT, refId)
		if err != nil {
			logger.Error("Failed to update email verification with sent time and resend email id", zap.Error(err))
		}
	}

	return c.String(http.StatusCreated, "Account created! Please verify your email to unlock all Konbini services.")
}

func handleVerifyEmail(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	routeErrorMessage := "Could not verify email. Please try again later."
	refId := c.QueryParam("code")

	ev, err := service.GetEmailVerification(refId)
	if err != nil {
		logger.Error("Error verifying email", zap.Error(err))
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	if ev.Status == service.EMAIL_STATUS_VERIFIED {
		logger.Error("Re-attempt to verify email.", zap.String("user_id", ev.UserId))
		return c.String(http.StatusBadRequest, "Verification code has been used before.")
	}

	// update the user entry that email has been verified
	logger.Info("Updating user entry to set email_verified...")
	tx, err := database.DB().Begin()
	if err != nil {
		logger.Error("Error begining transaction to verify email", zap.Error(err))
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}
	_, err = tx.Exec("UPDATE users SET email_verified = true WHERE id = $1;", ev.UserId)
	if err != nil {
		logger.Error("Error updating user entry to set email_verified to true", zap.Error(err))
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	// now we can update the email verification status because user entry has been updated
	logger.Info("Updating email verification status...")
	_, err = tx.Exec("UPDATE email_verifications SET status = $1, verified_at = $2 WHERE id = $3;", service.EMAIL_STATUS_VERIFIED, time.Now().In(time.UTC), ev.Id)
	if err != nil {
		// this error doesn't matter that much as long as the user entry has been updated
		logger.Error("Error updating email verification entry to set status to verified", zap.Error(err))
		database.Rollback(tx, c.Request().URL.Path)
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	err = tx.Commit()
	if err != nil {
		logger.Error("Error committing verify email changes", zap.Error(err))
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	logger.Info("Email verified")
	return c.String(http.StatusOK, "Email verified.")
}

func handleStartResetPassword(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	routeErrorMessage := "Could not start reset password process. Please try again later."

	// limit read to 1024 bytes or characters, more than enough for a single email
	logger.Info("Getting email from request body")
	bodyReader := io.LimitReader(c.Request().Body, 1024)
	reqBody, err := io.ReadAll(bodyReader)
	if err != nil {
		logger.Error("Failed to read request body", zap.Error(err))
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}
	email := string(reqBody)

	user, err := usermodel.GetByEmail(email)

	// generate random reset id
	logger.Info("Generating reset id...")
	resetId, err := gonanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 12)
	if err != nil {
		logger.Error("Failed to generate reset id", zap.Error(err))
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}
	logger.Info("Reset id", zap.String("reset_id", resetId))

	// expires in 1 hour
	expTime := time.Now().In(time.UTC).Add(time.Hour * 1)
	logger.Info("Expire time reset password", zap.Time("exp_time", expTime))

	// transaction
	logger.Info("Starting transaction to store reset password record")
	tx, err := database.DB().Begin()
	if err != nil {
		logger.Error("Failed to start transaction to store reset password record", zap.Error(err))
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	// check if there is another record of the user already, if there is, delete that one
	logger.Info("Checking if reset record exists for given user...")
	var recordId int64 = 0
	err = database.DB().QueryRow("SELECT id FROM users_passwords_resets WHERE user_id = $1;", user.Id).Scan(&recordId)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("Failed to check reset record", zap.Error(err))
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	if recordId != 0 {
		logger.Info("Updating existing password record...")
		_, err = tx.Exec("UPDATE users_passwords_resets SET reset_id = $1, expires_at = $2 WHERE id = $3;", resetId, expTime, recordId)
		if err != nil {
			database.Rollback(tx, "reset password")
			return c.String(http.StatusInternalServerError, routeErrorMessage)
		}
	} else {
		logger.Info("Saving new reset password record")
		_, err = tx.Exec("INSERT INTO users_passwords_resets (user_id, reset_id, expires_at) VALUES ($1, $2, $3);", user.Id, resetId, expTime)
		if err != nil {
			logger.Error("Failed to insert password reset link id into db", zap.Error(err))
			database.Rollback(tx, "reset password")
			return c.String(http.StatusInternalServerError, routeErrorMessage)
		}

	}

	err = tx.Commit()
	if err != nil {
		logger.Error("Failed to commit changes from transcation", zap.Error(err))
		database.Rollback(tx, "reset password")
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}
	logger.Info("Reset id saved", zap.String("reset_id", resetId))

	// send email with reset id
	if env.Values().APP_ENV == env.PRODUCTION {
		logger.Info("Generating email template...")
		var tpl bytes.Buffer
		err = templates.Render(&tpl, "reset-password.html", ResetPasswordEmailData{FirstName: user.FirstName, LastName: user.LastName, ResetId: resetId})
		logger.Info("Sending email with reset id for password reset...")
		_, err = service.SendEmail(env.Values().NOREPLY_EMAIL, email, "[Konbini] Reset Your Password", tpl.String())
		if err != nil {
			logger.Error("Failed to send password reset email", zap.Error(err))
			return c.String(http.StatusInternalServerError, routeErrorMessage)
		}
	}

	return c.String(http.StatusOK, fmt.Sprintf("Code: %s", resetId))
}

func handleFinishResetPassword(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	routeErrorMessage := "Could not reset password. Please try again later."
	// password reset codes are only 12 characters long
	reqBody := new(ResetPasswordEmailData)
	logger.Info("Reading password reset code from request body...")
	// bind the incoming request data
	if err := c.Bind(reqBody); err != nil {
		logger.Error("Error binding request body", zap.Error(err))
		return c.String(http.StatusBadRequest, "Invalid payload")
	}

	if err := c.Validate(reqBody); err != nil {
		logger.Error("Error validating request body", zap.Error(err))
		return err
	}

	logger.Info("Verifiying if code is valid", zap.String("reset_id", reqBody.ResetId))
	var (
		id        int64
		userId    string
		expiresAt time.Time
	)
	err := database.DB().QueryRow("SELECT id, user_id, expires_at FROM users_passwords_resets WHERE reset_id = $1;", reqBody.ResetId).Scan(&id, &userId, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("No row found for given reset code.", zap.String("code", reqBody.ResetId))
			return c.String(http.StatusBadRequest, "Invalid code or expired.")
		}
		logger.Error("Failed to get password reset record from db", zap.Error(err))
		return c.String(http.StatusInternalServerError, routeErrorMessage)
	}

	if time.Now().Before(expiresAt) {
		logger.Info("Starting transaction to perform password reset.")
		tx, err := database.DB().Begin()
		if err != nil {
			logger.Error("Failed to start transaction to perform password reset", zap.Error(err))
			return c.String(http.StatusInternalServerError, routeErrorMessage)
		}

		logger.Info("Deleting password reset record...")
		_, err = tx.Exec("DELETE FROM users_passwords_resets WHERE id = $1;", id)
		if err != nil {
			logger.Error("Failed to delete password reset record", zap.Error(err))
			database.Rollback(tx, "finish reset password")
			return c.String(http.StatusInternalServerError, routeErrorMessage)
		}

		logger.Info("Setting new password", zap.Int64("user_id", id))
		_, err = tx.Exec("UPDATE users SET password = crypt($1, gen_salt($2)) WHERE id = $3;", reqBody.Password, env.Values().PASS_ENCRYPT_ALGO, userId)
		if err != nil {
			logger.Error("Failed to update password", zap.Error(err))
			database.Rollback(tx, "finish reset password")
			return c.String(http.StatusInternalServerError, routeErrorMessage)
		}

		// TODO: add step to revoke all refresh tokens

		err = tx.Commit()
		if err != nil {
			logger.Error("Failed to commit password reset changes", zap.Error(err))
			database.Rollback(tx, "finish reset password")
			return c.String(http.StatusInternalServerError, routeErrorMessage)
		}
	} else {
		logger.Info("Expired password reset code", zap.String("code", reqBody.ResetId))
		logger.Info("Removing password reset code...")
		_, err = database.DB().Exec("DELETE FROM users_passwords_resets WHERE id = $1;", id)
		if err != nil {
			logger.Error("Failed to remove expired reset password code", zap.Error(err))
		}
		return c.String(http.StatusBadRequest, "Invalid code or expired.")
	}

	return c.String(http.StatusOK, "Password reset done. Please sign-in again on all your devices.")
}
