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

func VerifyEmail(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")
		if token == "" {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Missing token query parameter.",
			}
		}

		logger := middlewares.GetLogger(c)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		claims, err := services.ParseEmailToken(token)
		if err != nil {
			return err
		}

		queries := db.New(conn)

		emailToken, err := queries.GetEmailTokenById(ctx, claims.ID)

		// verify the token, this also verifies that the token hasn't expired yet.
		_, err = services.VerifyEmailToken(token, emailToken.TokenSalt)
		if err != nil {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Invalid token. Please request a new verification email.",
				InternalError: err,
			}
		}

		if claims.UserId != emailToken.UserID {
			return APIError{
				Code:           http.StatusBadRequest,
				PublicMessage:  "Invalid token. Please request a new verification email.",
				PrivateMessage: "The user id in the claims does not match the user id stored in the database.",
			}
		}

		tx, err := conn.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		queries = queries.WithTx(tx)
		ids, err := queries.DeleteAllEmailTokensByUserId(ctx, emailToken.UserID)
		if err != nil {
			return err
		}

		logger.Info().Strs("email_token_ids", ids).Msg("Deleted email tokens")

		err = queries.SetUserEmailVerifiedStatus(ctx, db.SetUserEmailVerifiedStatusParams{ID: emailToken.UserID, EmailVerified: true})
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

// LoginRequest represnets the request body for login route
type MagicLinkRequestRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

const magicLinkCodeLen = 6

func HandleMagicLinkRequest(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		body, err := middlewares.GetJsonBody[MagicLinkRequestRequest](c)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		queries := db.New(conn)

		exists, err := queries.ExistsUserWithEmail(ctx, body.Email)
		if err != nil {
			return err
		}

		if exists != 1 {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Invalid credentials.",
			}
		}

		user, err := queries.GetUserByEmail(ctx, body.Email)
		if err != nil {
			return err
		}

		matches, err := utils.ComparePasswordAndHash(body.Password, user.Password)
		if err != nil {
			return err
		}

		if !matches {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Invalid credentials.",
			}
		}

		// generate random digit sequence for the magic link
		digits, err := utils.RandomDigits(magicLinkCodeLen)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		exp := now.Add(time.Minute * 10)
		err = queries.CreateMagicLink(ctx, db.CreateMagicLinkParams{
			Token:     string(digits),
			UserID:    user.ID,
			CreatedAt: now.Format(time.RFC3339),
			ExpiresAt: exp.Format(time.RFC3339),
		})
		if err != nil {
			return err
		}

		// logger with request context
		logger := middlewares.GetLogger(c)
		go func(to, code, createdAt, expiresAt string, logger *zerolog.Logger) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			res, err := services.SendMagicLinkEmail(ctx, to, code, createdAt, expiresAt)
			if err != nil {
				logger.Error().Err(err).Str("to", to).Msg("Failed to send magic link email")
				return
			}
			logger.Info().Str("email_id", res.Id).Msg("Successfully sent magic link email")
		}(user.Email, string(digits), now.Format("2006/01/02 15:04PM")+" UTC", exp.Format("2006/01/02 15:04PM")+" UTC", logger)

		return c.NoContent(http.StatusOK)
	}
}

type MagicLinkVerifyRequest struct {
	Token string `json:"token" validate:"required,len=6"`
	Email string `json:"email" validate:"required,email"`
}

func HandleMagicLinkVerify(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

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
		defer conn.Close()
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

			conn, err := connector.Connect()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to connect to database when sending verification email")
				return
			}
			queries := db.New(conn)

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
