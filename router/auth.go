package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/juancwu/konbini/email"
	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// SetupAuthRouter is a helper function that will register all the auth routes to the RouterGroup.
func SetupAuthRouter(e RouterGroup) {
	e.POST("/auth/signup", handleSignup)
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
