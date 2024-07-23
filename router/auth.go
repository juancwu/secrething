package router

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/juancwu/konbini/email"
	"github.com/juancwu/konbini/jwt"
	"github.com/juancwu/konbini/middleware"
	"github.com/juancwu/konbini/store"
	"github.com/juancwu/konbini/views"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// SetupAuthRouter is a helper function that will register all the auth routes to the RouterGroup.
func SetupAuthRouter(e RouterGroup) {
	// sessions related routes
	e.POST("/auth/signup", handleSignup)
	e.POST("/auth/signin", handleSignin)
	e.PATCH("/auth/refresh", handleRefresh)

	// email related routes
	e.GET("/auth/email/verify", handleVerifyEmail)
	e.POST("/auth/email/resend", handleResendVerificationEmail)

	// account related routes
	e.GET("/auth/forgot/password", handleForgotPassword)
	e.GET("/auth/reset/password", handleResetPasswordForm)
	e.POST("/auth/reset/password", handleResetPassword)
	e.DELETE("/auth/account", handleDeleteAccount, middleware.Protect())
}

// handleSignup handles incoming signup requests
// This handler will create a new user and store it in the database.
func handleSignup(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	body := new(signupReqBody)

	log.Info().Msg("Binding signup request body.")
	err := c.Bind(body)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to bind signup request body.",
			Err:       err,
			RequestId: requestId,
		}
	}

	log.Info().Msg("Validating signup request body.")
	err = c.Validate(body)
	if err != nil {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Error when validating signup request body.",
			Err:       err,
			RequestId: requestId,
		}
	}

	log.Info().Msg("Checking for existing user with same email before creating a new user.")
	exists, err := store.ExistsUserWithEmail(body.Email)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Error when checking for existing user with email.",
			Err:       err,
			RequestId: requestId,
		}
	}

	if exists {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Existing user with the same email found. Abort user creation.",
			PublicMsg: "User with the given email already exists. If you forgot your password, please reset your password.",
			RequestId: requestId,
		}
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Creating new user.")
	user, err := store.NewUser(body.Email, body.Password, body.Name)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to create new user.",
			Err:       err,
			RequestId: requestId,
		}
	}
	log.Info().Str("email", user.Email).Str("user_id", user.Id).Msg("New user created.")

	// try to send email verification
	go sendVerificationEmail(requestId, user)

	return c.JSON(
		http.StatusCreated,
		map[string]string{
			"message": "Successfully signed up! Please check your email to verify it.",
		},
	)
}

// handleSignin handles incoming request to signin.
// An access token and a refresh token will be generated and send back to the client.
func handleSignin(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	body := new(signinReqBody)

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Binding signin request body.")
	err := c.Bind(body)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to bind signup request body.",
			Err:       err,
			RequestId: requestId,
		}
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Validating signin request body.")
	err = c.Validate(body)
	if err != nil {
		return apiError{
			Code: http.StatusBadRequest,
			Msg:  "Error when validating signin request body.",
			Err:  err,
		}
	}

	user, err := store.GetUserWithEmailAndPassword(body.Email, body.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusBadRequest,
				Msg:       "No user match with given email and password.",
				PublicMsg: "Invalid credentials. Please double check they are right.",
				Err:       err,
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to match user with email and password.",
			Err:       err,
			RequestId: requestId,
		}
	}

	tx, err := store.StartTx()
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to start transaction for token generation.",
			Err:       err,
			RequestId: requestId,
		}
	}

	at, err := jwt.GenerateAccessToken(tx, user.Id)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to generate access token.",
			Err:       err,
			RequestId: requestId,
		}
	}

	rt, err := jwt.GenerateRefreshToken(tx, user.Id)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to generate refresh token.",
			Err:       err,
			RequestId: requestId,
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback changes for the generated tokens when signing in.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to commit changes for the generated tokens when signing in.",
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, map[string]string{"access_token": at, "refresh_token": rt})
}

// handleVerifyEmail handles incoming request to verify an email of a user.
func handleVerifyEmail(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	code := c.QueryParam("code")
	if code == "" {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Invalid request. Missing code query parameter.",
			PublicMsg: "Invalid request. Missing code query parameter.",
			RequestId: requestId,
		}
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Get email verification based on code.")

	ev, err := store.GetEmailVerificationWithCode(code)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusBadRequest,
				Msg:       "No code found",
				PublicMsg: "Invalid code. Please get a new code.",
				Err:       err,
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to get email verification with code.",
			Err:       err,
			RequestId: requestId,
		}
	}

	now := time.Now()
	if now.After(ev.ExpiresAt) {
		// use go routine because the deletion of the expired email verification code is not essential to the request itself and prevents the client to wait longer than needed.
		go func() {
			// delete code so that it doesn't take up more space
			log.Info().Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code_id", ev.Id).Msg("Deleting expired code.")
			tx, err := store.StartTx()
			if err != nil {
				log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code_id", ev.Id).Msg("Failed to start transaction. Did not delete expired email verification code.")
			} else {
				_, err = ev.Delete(tx)
				if err != nil {
					log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code_id", ev.Id).Msg("Failed to delete email verification code.")
					err = tx.Rollback()
					if err != nil {
						log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code_id", ev.Id).Msg("Failed to rollback!")
					}
					return
				}
				err = tx.Commit()
				if err != nil {
					log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code_id", ev.Id).Msg("Failed to commit changes when deleting expired email verification code.")
				}
				log.Info().Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code_id", ev.Id).Msg("Expired code deleted.")
			}
		}()
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Email verification code expired.",
			PublicMsg: "Invalid code. Please get a new code.",
			RequestId: requestId,
		}
	}

	user, err := store.GetUserWithId(ev.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusBadRequest,
				Msg:       "Code has a user id that does not exists anymore. Check migrations if a proper cascade has been set.",
				PublicMsg: "Invalid code.",
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to get user",
			Err:       err,
			RequestId: requestId,
		}
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Start transaction to update user.")
	tx, err := store.StartTx()
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to start transaction to update user.",
			Err:       err,
			RequestId: requestId,
		}
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Setting user email verified to true.")
	res, err := user.SetEmailVerifiedTx(tx, true)
	if err != nil {
		store.Rollback(tx, requestId)
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to update user.",
			Err:       err,
			RequestId: requestId,
		}
	}
	_, err = res.RowsAffected()
	if err != nil {
		store.Rollback(tx, requestId)
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to update user.",
			RequestId: requestId,
		}
	}
	log.Info().Str(echo.HeaderXRequestID, requestId).Str("user_id", user.Id).Bool("email_verified", user.EmailVerified).Msg("User updated.")

	// now we have to delete the used email verification code
	log.Info().Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code", ev.Id).Msg("Delete used email verification code.")
	res, err = ev.Delete(tx)
	if err != nil {
		store.Rollback(tx, requestId)
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to delete used email verification code.",
			Err:       err,
			RequestId: requestId,
		}
	}
	_, err = res.RowsAffected()
	if err != nil {
		store.Rollback(tx, requestId)
		log.Error().Err(err).Int64("email_verification_code_id", ev.Id).Str(echo.HeaderXRequestID, requestId).Msg("Failed to delete used email verification code.")
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to delete used email verification code.",
			RequestId: requestId,
		}
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Committing changes in transaction.")
	err = tx.Commit()
	if err != nil {
		store.Rollback(tx, requestId)
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to commit changes in transaction.",
			Err:       err,
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, basicRespBody{Msg: "Successfully verifie email!", RequestId: requestId})
}

// handleResendVerificationEmail handles incoming request to send a new verification email.
func handleResendVerificationEmail(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	body := new(resendVerificationEmailReqBody)

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Binding resend verification email body.")
	err := c.Bind(body)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to bind signup request body.",
			Err:       err,
			RequestId: requestId,
		}
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Validating resend verification email body.")
	err = c.Validate(body)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Error when validating resend verificiation email body.",
			Err:       err,
			RequestId: requestId,
		}
	}

	// get user
	user, err := store.GetUserWithEmail(body.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusBadRequest,
				Msg:       "No user found.",
				PublicMsg: "No user found with the given email.",
				Err:       err,
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to get user with email",
			Err:       err,
			RequestId: requestId,
		}
	}

	if user.EmailVerified {
		err = errors.New("User's email has already been verified.")
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       err.Error(),
			PublicMsg: err.Error(),
			Err:       err,
			RequestId: requestId,
		}
	}

	// send the new email
	go sendVerificationEmail(requestId, user)

	return c.JSON(
		http.StatusOK,
		map[string]string{
			"message": "New verification email sent.",
		},
	)
}

// handleRefresh handles incoming requests for a new access token without using credentials but client must provide a valid refresh token.
func handleRefresh(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	if authHeader == "" {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       "Missing authorization header.",
			PublicMsg: http.StatusText(http.StatusUnauthorized),
			RequestId: requestId,
		}
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       "Invalid authorization header type.",
			PublicMsg: http.StatusText(http.StatusUnauthorized),
			RequestId: requestId,
		}
	}
	claims, err := jwt.VerifyRefreshToken(parts[1])
	if err != nil {
		return apiError{
			Code:      http.StatusUnauthorized,
			Err:       err,
			Msg:       "Failed to verify refresh token.",
			PublicMsg: http.StatusText(http.StatusUnauthorized),
			RequestId: requestId,
		}
	}
	user, err := store.GetUserWithId(claims.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusUnauthorized,
				Err:       err,
				Msg:       "Failed to get user.",
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get user.",
			RequestId: requestId,
		}
	}
	// create new access token
	tx, err := store.StartTx()
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to start transaction to generate access token using refresh token as authorization method.",
			RequestId: requestId,
		}
	}
	at, err := jwt.GenerateAccessToken(tx, user.Id)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to generate access token.",
			RequestId: requestId,
		}
	}
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollbakc changes when generating access token from refresh token.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to commit changes to save the newly genereated access token from refresh token.",
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, map[string]string{"access_token": at})
}

// Handle incoming request to start the password reset process.
// Users receive an email with a code and link to open in their browsers to finish the process.
// The link does not include the code itself but just the email. The email is used to put in the form that is
// then sent to the backend for further processing along with the code in their email.
func handleForgotPassword(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	userEmail := c.QueryParam("email")
	if userEmail == "" {
		return apiError{
			Code:      http.StatusBadRequest,
			PublicMsg: "Missing 'email' query parameter.",
			RequestId: requestId,
		}
	}

	user, err := store.GetUserWithEmail(userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusBadRequest,
				PublicMsg: "No user with given email.",
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to get user with email",
			Err:       err,
			RequestId: requestId,
		}
	}

	prc, err := store.NewOrUpdatePasswordResetCode(user.Id)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to create or update password reset code.",
			RequestId: requestId,
		}
	}

	template, err := email.RenderPasswordResetCodeEmail(user.Name, prc.Code, fmt.Sprintf("%s/auth/reset/password?email=%s", os.Getenv("SERVER_URL"), user.Email))
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to render password reset code email template.",
			RequestId: requestId,
		}
	}

	go func(template, requestId string) {
		res, err := email.Send("Reset Password - Konbini", os.Getenv("DONOTREPLY_EMAIL"), []string{user.Email}, template)
		if err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to send password reset code email.")
		}
		log.Info().Str(echo.HeaderXRequestID, requestId).Str("resend_email_id", res.Id).Msg("Password reset code email sent.")
	}(template, requestId)

	return writeJSON(http.StatusOK, c, map[string]string{"message": "You should receive an email with a link to reset your password."})
}

// Handles requests to reset a user's password.
// A code of length 6 is required and gotten from the forgot passwor route.
func handleResetPassword(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	email := c.FormValue("email")
	code := c.FormValue("code")
	password := c.FormValue("password")

	if email == "" || code == "" || password == "" {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Missing code, email, or password form values.",
			PublicMsg: "Missing code, email, or password form values.",
			RequestId: requestId,
		}
	}

	// TODO: add password format validation

	user, err := store.GetUserWithEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusBadRequest,
				Msg:       fmt.Sprintf("No user found with given email. Email: %s", email),
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get user",
			RequestId: requestId,
		}
	}

	prc, err := store.GetPasswordResetCodeByUserId(user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusBadRequest,
				Msg:       "No password reset code found",
				PublicMsg: "Invalid code.",
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to get password reset code",
			Err:       err,
			RequestId: requestId,
		}
	}

	if time.Now().After(prc.ExpiresAt) {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Password reset code is expired.",
			PublicMsg: "Invalid code.",
			RequestId: requestId,
		}
	}

	if prc.Code != code {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Password reset codes do not match.",
			PublicMsg: "Invalid code.",
			RequestId: requestId,
		}
	}

	// update the user's password
	user.Password = password
	tx, err := store.StartTx()
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to start transaction to update user password.",
			Err:       err,
			RequestId: requestId,
		}
	}
	if _, err := user.Update(tx); err != nil {
		rollbakErr := tx.Rollback()
		if rollbakErr != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to update user password",
			Err:       err,
			RequestId: requestId,
		}
	}

	if _, err := prc.Delete(tx); err != nil {
		rollbakErr := tx.Rollback()
		if rollbakErr != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to delete password reset code after usage.",
			Err:       err,
			RequestId: requestId,
		}
	}

	// delete the tokens associated to user
	deleteResult, err := store.DeleteTokensOwnedByUser(tx, user.Id)
	if err != nil {
		rollbakErr := tx.Rollback()
		if rollbakErr != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to delete all tokens owned by user.",
			Err:       err,
			RequestId: requestId,
		}
	}
	if n, err := deleteResult.RowsAffected(); err != nil {
		log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to get how many rows from auth_tokens where deleted.")
	} else {
		log.Info().Int64("delete_count", n).Str(echo.HeaderXRequestID, requestId).Msg("Rows deleted from auth_tokens.")
	}

	if err := tx.Commit(); err != nil {
		rollbakErr := tx.Rollback()
		if rollbakErr != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to commit changes",
			Err:       err,
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, map[string]string{"message": "Password reset successful.", "request_id": requestId})
}

// Only handles serving the reset password form view.
func handleResetPasswordForm(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	email := c.QueryParam("email")

	var html bytes.Buffer
	if err := views.ResetPasswordForm(fmt.Sprintf("%s/auth/reset/password", os.Getenv("SERVER_URL")), email).Render(context.Background(), &html); err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			RequestId: requestId,
		}
	}

	return c.HTML(http.StatusOK, html.String())
}

// Handle requests to delete user account. Everything related to the account will be deleted.
// This includes all prepared bentos.
func handleDeleteAccount(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	claims, err := middleware.GetJwtClaimsFromContext(c)
	if err != nil {
		if errors.Is(err, middleware.ErrNoJwtClaims) {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get jwt claims from context to delete account.",
			RequestId: requestId,
		}
	}

	user, err := store.GetUserWithId(claims.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusNotFound,
				Err:       err,
				Msg:       "No user found to delete.",
				PublicMsg: "No user found to delete.",
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get user with id to delete account.",
			RequestId: requestId,
		}
	}

	tx, err := store.StartTx()
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to start transaction to delete user.",
			RequestId: requestId,
		}
	}

	res, err := user.Delete(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to delete user",
			RequestId: requestId,
		}
	}

	n, err := res.RowsAffected()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get affected rows count after perfroming DELETE on user.",
			RequestId: requestId,
		}
	}

	if n < 1 {
		return apiError{
			Code:      http.StatusNotModified,
			Msg:       "Failed to delete account since less than 1 users where deleted.",
			PublicMsg: "Account not deleted.",
			RequestId: requestId,
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to commit changes to delete account.",
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, basicRespBody{Msg: "Account deleted. All data related to account pruned. This include all bentos previously created.", RequestId: requestId})
}
