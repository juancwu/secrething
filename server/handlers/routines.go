package handlers

import (
	"context"
	"konbini/server/services"
	"time"

	"github.com/rs/zerolog"
)

// sendVerificationEmail is a helper function that sends a verification email to the given user email.
// The function is intended to be used as a go routine and it will log any error with the provided logger.
// The function will create a new email token and store the token in memory cache using storeEmailTokenInCache.
// The stored email token can later to retrieved by getEmailTokenFromCache using the id.
func sendVerificationEmail(userId string, userEmail string, logger *zerolog.Logger) {
	// sending an email shouldn't take more than 1 minute
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	token, err := services.NewEmailToken(userId)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create email token.")
		return
	}

	err = storeEmailTokenInCache(token)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to store email token in memory cache.")
		return
	}

	tokenStr, err := token.Package()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to package email token.")
		return
	}

	res, err := services.SendVerificationEmail(ctx, userEmail, tokenStr)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to send verification email")
		return
	}

	logger.Info().
		Str("email_id", res.Id).
		Str("user_id", userId).
		Str("user_email", userEmail).
		Msg("Successfully sent verification email")
}
