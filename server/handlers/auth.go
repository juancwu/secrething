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
)

// RegisterRequest represents the request body for register route.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=32"`
	NickName string `json:"nickname" validate:"required,min=3,max=32"`
}

// Register is a handler function that registers a user for Konbini.
func Register(connector *db.DBConnector) echo.HandlerFunc {
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

		// userId at
		now := time.Now().UTC().Format(time.RFC3339)

		userId, err := queries.CreateUser(ctx, db.CreateUserParams{
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

		logger.Info().Str("user_id", userId).Msg("New user registered.")
		go sendVerificationEmail(userId, body.Email, logger)

		return c.NoContent(http.StatusCreated)
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Login(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := middlewares.GetJsonBody[LoginRequest](c)
		if err != nil {
			return err
		}

		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		queries := db.New(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

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
				PublicMessage: "Invalid credentials. Please try again.",
			}
		}

		return nil
	}
}

func VerifyEmail(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")
		if token == "" {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Missing token query parameter.",
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		id, err := services.ExtractEmailTokenId(token)
		if err != nil {
			return APIError{
				Code:           http.StatusBadRequest,
				PublicMessage:  "Invalid.",
				PrivateMessage: "Failed to extract email token id from token string.",
				InternalError:  err,
			}
		}

		emailToken, err := getEmailTokenFromCache(id)
		if err != nil {
			return err
		}

		if time.Now().UTC().After(emailToken.ExpiresAt) {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Expired.",
			}
		}

		userId := emailToken.UserId
		queries := db.New(conn)
		err = queries.SetUserEmailVerifiedStatus(ctx, db.SetUserEmailVerifiedStatusParams{ID: userId, EmailVerified: true})
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

func SetupTOTP(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func VerifyTOTP(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
