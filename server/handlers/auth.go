package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"konbini/server/config"
	"konbini/server/db"
	"konbini/server/middlewares"
	"konbini/server/services"
	"konbini/server/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

		if claims.Subject != emailToken.UserID {
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

// TODO: update the client validation to custom magic link client validation

// LoginRequest represnets the request body for login route
type MagicLinkRequestRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Client   string `json:"client" validate:"required,oneof=tenin konbi"`
	// State is a string that the client sends and will be included in the ramaining
	// steps of the magic link authentication process. The client should implement
	// its own method to verify the integrity of the state string.
	State string `json:"state" validate:"required,max=1024"`
}

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

		now := time.Now().UTC()
		exp := now.Add(time.Minute * 10)

		magicLinkId, err := queries.CreateMagicLink(ctx, db.CreateMagicLinkParams{
			UserID:    user.ID,
			State:     body.State,
			CreatedAt: now.Format(time.RFC3339),
			ExpiresAt: exp.Format(time.RFC3339),
		})
		if err != nil {
			return err
		}

		// logger with request context
		logger := middlewares.GetLogger(c)
		go sendMagicLinkRoutine(
			magicLinkRoutineParams{
				LinkId:    magicLinkId,
				Email:     user.Email,
				UserId:    user.ID,
				CreatedAt: now.Format(time.RFC3339),
				ExpiresAt: exp.Format(time.RFC3339),
				Client:    body.Client,
				State:     body.State,
				Logger:    logger,
			},
		)

		return c.JSON(http.StatusOK, map[string]string{"magic_link_id": magicLinkId})
	}
}

type magicLinkRoutineParams struct {
	LinkId    string
	Email     string
	UserId    string
	CreatedAt string
	ExpiresAt string
	Client    string
	State     string
	Logger    *zerolog.Logger
}

func sendMagicLinkRoutine(params magicLinkRoutineParams) {
	logger := params.Logger
	linkId := params.LinkId
	userId := params.UserId
	email := params.Email
	createdAt := params.CreatedAt
	expiresAt := params.ExpiresAt
	client := params.Client
	state := params.State

	c, err := config.Global()
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("Failed to get server configuration.")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// all query parameters are to be encrypted to avoid revealing any of this information if
	// the email ends up being seen by some unwanted entity.
	magicLinkContext := userId + linkId
	encryptedMagicLinkContext, err := utils.EncryptAES([]byte(magicLinkContext), c.GetAesKey())
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("Failed to encrypt user id and code.")
		return
	}

	encryptedClient, err := utils.EncryptAES([]byte(client), c.GetAesKey())
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("Failed to encrypt client.")
		return
	}

	encryptedState, err := utils.EncryptAES([]byte(state), c.GetAesKey())
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("Failed to encrypt state.")
		return
	}

	magicUrl := fmt.Sprintf(
		"%s/api/v1/auth/magic/verify?token=%s&client=%s&state=%s",
		c.GetBackendUrl(),
		base64.URLEncoding.EncodeToString(encryptedMagicLinkContext),
		base64.URLEncoding.EncodeToString(encryptedClient),
		base64.URLEncoding.EncodeToString(encryptedState),
	)
	res, err := services.SendMagicLinkEmail(ctx, email, magicUrl, createdAt, expiresAt)
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("Failed to send magic link email")
		return
	}
	logger.Info().Str("email_id", res.Id).Msg("Successfully sent magic link email")
}

func HandleMagicLinkVerify(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg, err := config.Global()
		if err != nil {
			return err
		}
		cache, err := cfg.Cache()
		if err != nil {
			return err
		}

		b64Token := c.QueryParam("token")
		b64Client := c.QueryParam("client")
		b64State := c.QueryParam("state")

		if len(b64Token) < 1 || len(b64Client) < 1 || len(b64State) < 1 {
			return APIError{
				Code:           http.StatusBadRequest,
				PublicMessage:  "Invalid magic link",
				PrivateMessage: "Token, client, or state query paramters too short",
			}
		}

		decoded, err := base64.URLEncoding.DecodeString(b64Token)
		if err != nil {
			return err
		}
		magicLinkContext, err := utils.DecryptAES(decoded, cfg.GetAesKey())
		if err != nil {
			return err
		}
		if len(magicLinkContext) < 72 {
			return APIError{
				Code:           http.StatusBadRequest,
				PublicMessage:  "Invalid magic link",
				PrivateMessage: "Magic link context is too short",
			}
		}

		decoded, err = base64.URLEncoding.DecodeString(b64Client)
		if err != nil {
			return err
		}
		client, err := utils.DecryptAES(decoded, cfg.GetAesKey())
		if err != nil {
			return err
		}
		if string(client) != "tenin" && string(client) != "konbi" {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Invalid client value in magic link",
			}
		}

		decoded, err = base64.URLEncoding.DecodeString(b64State)
		if err != nil {
			return err
		}
		state, err := utils.DecryptAES(decoded, cfg.GetAesKey())
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		// connect to database
		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		queries := db.New(conn)

		// grab the magic link in database
		userId := magicLinkContext[:36]
		magicLinkId := magicLinkContext[36:]
		magicLink, err := queries.GetMagicLink(ctx, db.GetMagicLinkParams{ID: string(magicLinkId), UserID: string(userId), State: string(state)})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "Invalid magic link",
					InternalError: err,
				}
			}
			return err
		}

		// NOTE: skip check on state since state should be the same otherwise the query wouldn't return anything

		// check magic link expiration
		exp, err := time.Parse(time.RFC3339, magicLink.ExpiresAt)
		if err != nil {
			return err
		}

		if time.Now().UTC().After(exp) {
			// make sure to remove old magic link
			go removeInvalidMagicLink(connector, magicLink.ID, magicLink.UserID, magicLink.State)
			return APIError{
				Code:          http.StatusBadGateway,
				PublicMessage: "Expired magic link",
			}
		}

		// at this point, the magic link is valid and we should remove it since it has been used
		go removeInvalidMagicLink(connector, magicLink.ID, magicLink.UserID, magicLink.State)

		hash := sha256.New()
		hash.Write(magicLinkId)
		hash.Write(client)
		hash.Write(state)
		key := hash.Sum(nil)
		log.Info().Str("key", string(key)).Str("id", string(magicLinkId)).Str("client", string(client)).Str("state", string(state)).Msg("from verify")
		// store that the magic link has been validated
		// give clients 1 minute to get their access token
		cache.Set(string(key), true, time.Minute)

		return c.String(http.StatusOK, "success")
	}
}

func removeInvalidMagicLink(connector *db.DBConnector, id, userId, state string) {
	conn, err := connector.Connect()
	if err != nil {
		return
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	db.New(conn).RemoveMagicLink(ctx, db.RemoveMagicLinkParams{ID: id, UserID: userId, State: state})
}

type MagicLinkStatusRequest struct {
	Id     string `json:"id" validate:"required,uuid"`
	Client string `json:"client" validate:"required,oneof=tenin konbi"`
	State  string `json:"state" validate:"required,max=1024"`
}

// Handle requests to check if a magic link status has been verified or not
func HandleMagicLinkStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg, err := config.Global()
		if err != nil {
			return err
		}
		cache, err := cfg.Cache()
		if err != nil {
			return err
		}
		body, err := middlewares.GetJsonBody[MagicLinkStatusRequest](c)
		if err != nil {
			return err
		}
		hash := sha256.New()
		hash.Write([]byte(body.Id))
		hash.Write([]byte(body.Client))
		hash.Write([]byte(body.State))
		key := hash.Sum(nil)
		log.Info().Str("key", string(key)).Str("id", body.Id).Str("client", body.Client).Str("state", body.State).Msg("from status")
		_, found := cache.Get(string(key))
		if found {
			// one time use only
			cache.Delete(string(key))
			return c.NoContent(http.StatusOK)
		}
		return c.NoContent(http.StatusBadRequest)
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
