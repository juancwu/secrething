package handlers

import (
	"context"
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

// LoginRequest represnets the request body for login route
type MagicLinkRequestRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	Client      string `json:"client" validate:"required,oneof=tenin konbi"`
	RedirectUri string `json:"redirect_uri" validate:"required,http_url"`
	State       string `json:"state" validate:"required,max=1024"`
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
		go sendMagicLinkRoutine(
			magicLinkRoutineParams{
				Email:       user.Email,
				UserId:      user.ID,
				Code:        string(digits),
				CreatedAt:   now.Format(time.RFC3339),
				ExpiresAt:   exp.Format(time.RFC3339),
				Client:      body.Client,
				RedirectUri: body.RedirectUri,
				State:       body.State,
				Logger:      logger,
			},
		)

		return c.NoContent(http.StatusOK)
	}
}

type magicLinkRoutineParams struct {
	Email       string
	UserId      string
	Code        string
	CreatedAt   string
	ExpiresAt   string
	Client      string
	RedirectUri string
	State       string
	Logger      *zerolog.Logger
}

func sendMagicLinkRoutine(params magicLinkRoutineParams) {
	logger := params.Logger
	userId := params.UserId
	email := params.Email
	code := params.Code
	createdAt := params.CreatedAt
	expiresAt := params.ExpiresAt
	client := params.Client
	redirectUri := params.RedirectUri
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
	userContext := userId + code
	encryptedUserContext, err := utils.EncryptAES([]byte(userContext), c.GetAesKey())
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("Failed to encrypt user id and code.")
		return
	}

	encryptedRedirectUri, err := utils.EncryptAES([]byte(redirectUri), c.GetAesKey())
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("Failed to encrypt redirect uri.")
		return
	}

	encryptedClient, err := utils.EncryptAES([]byte(client), c.GetAesKey())
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("Failed to encrypt client.")
		return
	}

	magicUrl := fmt.Sprintf(
		"%s/api/v1/auth/magic/verify?token=%s&redirect_uri=%s&client=%s&state=%s",
		c.GetBackendUrl(),
		base64.URLEncoding.EncodeToString(encryptedUserContext),
		base64.URLEncoding.EncodeToString(encryptedRedirectUri),
		base64.URLEncoding.EncodeToString(encryptedClient),
		base64.URLEncoding.EncodeToString([]byte(state)),
	)
	res, err := services.SendMagicLinkEmail(ctx, email, magicUrl, createdAt, expiresAt)
	if err != nil {
		logger.Error().Err(err).Str("to", email).Msg("Failed to send magic link email")
		return
	}
	logger.Info().Str("email_id", res.Id).Msg("Successfully sent magic link email")
}

// magicLinkVerification is used primarily to verify the query parameters for the magic lin.
type magicLinkVerification struct {
	Client      string `validate:"required,oneof=tenin konin"`
	RedirectUri string `validate:"required,http_url"`
}

func HandleMagicLinkVerify(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get all the required query parameters
		b64Token := c.QueryParam("token")
		if b64Token == "" {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Missing token in magic link.",
			}
		}

		cfg, err := config.Global()
		if err != nil {
			return err
		}

		// make sure that all parameters are valid before querying the database
		redirectUri := c.QueryParam("redirect_uri")
		decodedRedirectUriLen := base64.URLEncoding.DecodedLen(len(redirectUri))
		decodedRedirectUri := make([]byte, decodedRedirectUriLen)
		_, err = base64.URLEncoding.Decode(decodedRedirectUri, []byte(redirectUri))
		if err != nil {
			return err
		}
		decryptedRedirectUri, err := utils.DecryptAES(decodedRedirectUri, cfg.GetAesKey())
		if err != nil {
			return err
		}
		if len(decryptedRedirectUri) == 0 {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Missing redirect_uri in magic link.",
			}
		}
		redirectUri = string(decryptedRedirectUri)

		client := c.QueryParam("client")
		decodedClientLen := base64.URLEncoding.DecodedLen(len(client))
		decodedClient := make([]byte, decodedClientLen)
		_, err = base64.URLEncoding.Decode(decodedClient, []byte(client))
		if err != nil {
			return err
		}
		decryptedClient, err := utils.DecryptAES(decodedClient, cfg.GetAesKey())
		if err != nil {
			return err
		}
		if len(decryptedClient) == 0 {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Missing client in magic link.",
			}
		}
		client = string(decryptedClient)

		// make sure they are valid client and redirect url
		if err := c.Validate(&magicLinkVerification{Client: client, RedirectUri: redirectUri}); err != nil {
			return APIError{
				Code:           http.StatusBadRequest,
				PublicMessage:  "Invalid magic link.",
				PrivateMessage: "Cient and/or redirect uir is/are invalid.",
				InternalError:  err,
			}
		}

		// decrypt the token and get the context
		decodedTokenLen := base64.URLEncoding.DecodedLen(len(b64Token))
		decodedToken := make([]byte, decodedTokenLen)
		_, err = base64.URLEncoding.Decode(decodedToken, []byte(b64Token))
		if err != nil {
			return err
		}
		decryptedToken, err := utils.DecryptAES(decodedToken, cfg.GetAesKey())
		if err != nil {
			return err
		}
		token := string(decryptedToken)
		// uuid v4 has 36 characters and the code has 6
		// 36 + 6 = 42
		if len(token) != 42 {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Invalid token.",
			}
		}
		userId := string(token[:36])
		code := string(token[36:])

		// get the magic link from the database
		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()
		queries := db.New(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		magicLink, err := queries.GetMagicLink(ctx, db.GetMagicLinkParams{
			Token:  code,
			UserID: userId,
		})
		if err != nil {
			return err
		}

		// check for expiration
		now := time.Now().UTC()
		exp, err := time.Parse(time.RFC3339, magicLink.ExpiresAt)
		if err != nil {
			return err
		}
		if now.After(exp.UTC()) {
			// remove the expired magic link in a go routine to not block
			go removeInvalidMagicLink(connector, magicLink.Token, magicLink.UserID)
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Expired magic link.",
			}
		}

		// check if code matches
		if magicLink.Token != code {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Invalid magic link.",
			}
		}

		// remove the magic link in go routine to not block
		go removeInvalidMagicLink(connector, magicLink.Token, magicLink.UserID)

		// create a new session in the database
		// this helps keep track how many devices the user has signed in
		// and also allow users to control which session to sign out.

		// Each session has its own token salt making revoking individual sessions possible
		salt, err := utils.RandomBytes(16)
		if err != nil {
			return err
		}

		// TODO: add code to determine if the same machine is login in again
		// TODO: Missing information gathering to differentiate sessions for viewing purposes

		// create a new session
		tokenId, err := queries.CreateSession(ctx, db.CreateSessionParams{
			TokenSalt:    salt,
			UserID:       magicLink.UserID,
			Ip:           sql.NullString{String: c.RealIP(), Valid: true},
			LastActivity: now.Format(time.RFC3339),
		})
		if err != nil {
			return err
		}

		// generate new user token
		exp = now.Add(7 * 24 * time.Hour).UTC() // a week
		token, err = services.NewUserToken(
			tokenId,
			magicLink.UserID,
			"okyaku-sama",
			"authentication",
			salt,
			exp,
		)
		if err != nil {
			return err
		}

		redirectUrl := fmt.Sprintf("%s?token=%s&client=%s", redirectUri, token, client)

		return c.Redirect(http.StatusPermanentRedirect, redirectUrl)
	}
}

func removeInvalidMagicLink(connector *db.DBConnector, token, userId string) {
	conn, err := connector.Connect()
	if err != nil {
		return
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	db.New(conn).RemoveMagicLink(ctx, db.RemoveMagicLinkParams{Token: token, UserID: userId})
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
