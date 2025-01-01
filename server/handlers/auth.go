package handlers

import (
	"context"
	"errors"
	"konbini/server/db"
	"konbini/server/middlewares"
	"konbini/server/services"
	"konbini/server/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// RegisterRequest represents the request body for register route.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=32"`
	NickName string `json:"nickname" validate:"required,min=3,max=32"`
}

// HandleRegister is a handler function that registers a user for Konbini.
func HandleRegister(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		queries := db.New(conn)
		body, ok := c.Get(middlewares.JSON_BODY_KEY).(*RegisterRequest)
		if !ok {
			return errors.New("Failed to get JSON body from context.")
		}

		logger := middlewares.GetLogger(c)

		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
		defer cancel()

		exists, err := queries.ExistsUserWithEmail(ctx, body.Email)
		if err != nil {
			return err
		}

		if exists == 1 {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Email has already been taken.",
			}
		}

		// create password hash
		hash, err := utils.GeneratePasswordHash(body.Password)
		if err != nil {
			return err
		}

		// random jwt salt
		salt, err := utils.RandomBytes(32)
		if err != nil {
			return err
		}

		// partialUserInformation at
		now := time.Now().UTC().Format(time.RFC3339)

		partialUserInformation, err := queries.CreateUser(ctx, db.CreateUserParams{
			Email:     body.Email,
			Password:  hash,
			Nickname:  body.NickName,
			TokenSalt: salt,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return err
		}

		logger.Info().Str("user_id", partialUserInformation.ID).Msg("New user partialUserInformation.")

		go func(userId string, userEmail string, logger *zerolog.Logger) {
			// sending an email shouldn't take more than 1 minute
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			salt, err := utils.RandomBytes(16)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to generate random email token key when sending verification email")
				return
			}

			now := time.Now().UTC()
			exp := now.Add(time.Minute * 10).UTC()
			createEmailTokenParams := db.CreateEmailTokenParams{
				UserID:    userId,
				TokenSalt: salt,
				CreatedAt: now.Format(time.RFC3339),
				ExpiresAt: exp.Format(time.RFC3339),
			}
			id, err := queries.CreateEmailToken(ctx, createEmailTokenParams)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create email token in database when sending verification email")
				return
			}

			token, err := services.NewEmailToken(id, userId, salt, exp)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to generate email jwt when sending verification email")
				return
			}

			userEmail = "delivered@resend.dev"
			res, err := services.SendVerificationEmail(ctx, userEmail, token)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to send verification email")
				return
			}

			logger.Info().
				Str("email_id", res.Id).
				Str("user_id", userId).
				Str("user_email", userEmail).
				Msg("Successfully sent verification email")

		}(partialUserInformation.ID, body.Email, logger)

		return c.NoContent(http.StatusCreated)
	}
}
