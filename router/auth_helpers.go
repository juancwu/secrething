package router

import (
	"fmt"
	"os"

	"github.com/juancwu/konbini/email"
	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// sendVerificationEmail is a helper method to send a new verification email.
// This method will take care of deleting any previous codes and create a new one.
func sendVerificationEmail(requestId string, user *store.User) {
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
}
