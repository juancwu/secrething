package handlers

import (
	"context"
	"fmt"
	"konbini/server/memcache"
	"konbini/server/services"
	"time"

	"github.com/rs/zerolog"
)

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

func storeEmailTokenInCache(token *services.EmailToken) error {
	cache := memcache.Cache()
	return cache.Add("email_token_"+token.Id, token, time.Minute*10)
}

func getEmailTokenFromCache(id string) (*services.EmailToken, error) {
	cache := memcache.Cache()
	k, exp, found := cache.GetWithExpiration("email_token_" + id)
	if !found {
		return nil, fmt.Errorf("No email token found with id: %s", id)
	}
	if time.Now().UTC().After(exp) {
		return nil, fmt.Errorf("Email token cache expired. ID: %s", id)
	}
	token, ok := k.(*services.EmailToken)
	if !ok {
		return nil, fmt.Errorf("Invalid email token type.")
	}
	return token, nil
}
