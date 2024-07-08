package router

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/juancwu/konbini/email"
	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// SetupAuthRouter is a helper function that will register all the auth routes to the RouterGroup.
func SetupAuthRouter(e RouterGroup) {
	e.POST("/auth/signup", handleSignup)
	e.GET("/auth/email/verify", handleVerifyEmail)
}

// handleSignup handles incoming signup requests
// This handler will create a new user and store it in the database.
func handleSignup(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	body := new(signupReqBody)

	log.Info().Msg("Binding signup request body.")
	err := c.Bind(body)
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to bind signup request body.")
		return err
	}

	log.Info().Msg("Validating signup request body.")
	err = c.Validate(body)
	if err != nil {
		c.Set(err_msg_logger_key, "Error when validating signup request body.")
		return err
	}

	log.Info().Msg("Checking for existing user with same email before creating a new user.")
	exists, err := store.ExistsUserWithEmail(body.Email)
	if err != nil {
		c.Set(err_msg_logger_key, "Error when checking for existing user with email.")
		return err
	}

	if exists {
		log.Info().Msg("Existing user with the same email found. Abort user creation.")
		c.Set(err_msg_logger_key, "Abort new user creation due to duplication.")
		c.Set(public_err_msg_key, "User with the given email already exists. If you forgot your password, please reset your password.")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Creating new user.")
	user, err := store.NewUser(body.Email, body.Password, body.Name)
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to create new user.")
		return err
	}
	log.Info().Str("email", user.Email).Str("user_id", user.Id).Msg("New user created.")

	// try to send email verification
	go func() {
		log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Creating email verification.")
		_, err := store.DeleteEmailVerificationWithUserId(user.Id)
		if err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to delete email verification code with user id. This is a step to prevent unique constraint issues when inserting a new row.")
		}
		// try to create a new email verification anyways
		ev, err := store.NewEmailVerification(user.Id)
		if err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to create email verification.")
		} else {
			url := fmt.Sprintf("%s/auth/email/verify?code=%s", os.Getenv("SERVER_URL"), ev.Code)
			html, err := email.RenderVerifiationEmail(user.Name, url)
			if err != nil {
				log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to render email verification html.")
			} else {
				sent, err := email.Send("Verify Your Email", os.Getenv("DONOTREPLY_EMAIL"), []string{user.Email}, html)
				if err != nil {
					log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to send verification email on user created.")
					return
				}
				log.Info().Str("email", user.Email).Str("resend_id", sent.Id).Str(echo.HeaderXRequestID, requestId).Msg("Verification email sent.")
			}
		}
	}()

	return c.JSON(
		http.StatusCreated,
		map[string]string{
			"message": "Successfully signed up! Please check your email to verify it.",
		},
	)
}

// handleVerifyEmail handles incoming request to verify an email of a user.
func handleVerifyEmail(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	code := c.QueryParam("code")
	if code == "" {
		c.Set(public_err_msg_key, "Invalid request. Missing code query parameter.")
		c.Set(err_msg_logger_key, "Invalid request. Missing code query parameter.")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Get email verification based on code.")

	ev, err := store.GetEmailVerificationWithCode(code)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Set(public_err_msg_key, "Invalid code. Please get a new code.")
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		c.Set(err_msg_logger_key, "Failed to get email verification with code.")
		return err
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
		c.Set(public_err_msg_key, "Invalid code. Please get a new code.")
		c.Set(err_msg_logger_key, "Email verification code expired.")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	user, err := store.GetUserWithId(ev.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Set(public_err_msg_key, "Invalid code.")
			c.Set(err_msg_logger_key, "Code has a user id that does not exists anymore. Check migrations if a proper cascade has been set.")
			log.Error().Err(errors.New("Code has invalid user id.")).Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code_id", ev.Id).Msg("Invalid user id in email verification record.")
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		return err
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Start transaction to update user.")
	tx, err := store.StartTx()
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to start transaction to update user.")
		return err
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Setting user email verified to true.")
	user.EmailVerified = true
	res, err := user.Update(tx)
	if err != nil {
		log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Str("user_id", user.Id).Msg("Failed to update user.")
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error().Err(rollbackErr).Str(echo.HeaderXRequestID, requestId).Str("user_id", user.Id).Msg("Failed to rollback after failing to update user.")
		}
		c.Set(err_msg_logger_key, "Failed to update user.")
		return err
	}
	n, _ := res.RowsAffected()
	if n < 1 {
		err = errors.New("Failed to update user model.")
		log.Error().Err(err).Str("user_id", user.Id).Str(echo.HeaderXRequestID, requestId).Msg("Failed to update user.")
		return err
	} else if n > 1 {
		log.Warn().Str(echo.HeaderXRequestID, requestId).Str("user_id", user.Id).Msgf("Multiple users where updated when trying to update one user. Count: %d. Will rollback.", n)
		err = tx.Rollback()
		if err != nil {
			c.Set(err_msg_logger_key, "Failed to rollback on multiple users updated on email verification process.")
			return err
		}
	}
	log.Info().Str(echo.HeaderXRequestID, requestId).Str("user_id", user.Id).Bool("email_verified", user.EmailVerified).Msg("User updated.")

	// now we have to delete the used email verification code
	log.Info().Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code", ev.Id).Msg("Delete used email verification code.")
	res, err = ev.Delete(tx)
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to delete used email verification code.")
		return err
	}
	n, _ = res.RowsAffected()
	if n < 1 {
		err = errors.New("Failed to delete used email verification code.")
		log.Error().Err(err).Int64("email_verification_code_id", ev.Id).Str(echo.HeaderXRequestID, requestId).Msg("Failed to delete used email verification code.")
		return err
	} else if n > 1 {
		log.Warn().Str(echo.HeaderXRequestID, requestId).Int64("email_verification_code_id", ev.Id).Msgf("Multiple email verifications where deleted when trying to delete one. Count: %d. Will rollback.", n)
		err = tx.Rollback()
		if err != nil {
			c.Set(err_msg_logger_key, "Failed to rollback on multiple deleted email verifications on email verification process.")
			return err
		}
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Committing changes in transaction.")
	err = tx.Commit()
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to commit changes in transaction.")
		return err
	}

	return c.JSON(
		http.StatusOK,
		map[string]string{
			"message": "Successfully verified email.",
		},
	)
}
