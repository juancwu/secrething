package handlers

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"konbini/server/db"
	"konbini/server/memcache"
	"konbini/server/middlewares"
	"konbini/server/services"
	"konbini/server/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
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

		// userId at
		now := time.Now().UTC().Format(time.RFC3339Nano)

		userId, err := queries.CreateUser(ctx, db.CreateUserParams{
			Email:     body.Email,
			Password:  hash,
			Nickname:  body.NickName,
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
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required"`
	TOTPCode *string `json:"totp_code,omitempty" validate:"omitnil,omitempty,required,len=6"`
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

		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Minute)
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

		var tokType services.TokenType
		if !user.TotpSecret.Valid || !user.EmailVerified {
			tokType = services.PARTIAL_USER_TOKEN_TYPE
		} else if user.TotpLocked && user.TotpSecret.Valid {
			if body.TOTPCode == nil {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "User has TOTP setup, code is required. Make a new login request with the totp_code field in the body.",
				}
			}
			if !totp.Validate(*body.TOTPCode, user.TotpSecret.String) {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "Invalid TOTP code.",
				}
			}
		} else {
			tokType = services.FULL_USER_TOKEN_TYPE
		}

		now := time.Now().UTC()
		exp := now.Add(time.Hour * 24 * 7)
		var authToken *services.AuthToken
		dbJwt, err := queries.NewJWT(ctx, db.NewJWTParams{
			UserID:    user.ID,
			TokenType: tokType.String(),
			CreatedAt: now.Format(time.RFC3339Nano),
			ExpiresAt: exp.Format(time.RFC3339Nano),
		})
		authToken, err = services.NewAuthToken(dbJwt.ID, user.ID, tokType, exp)
		if err != nil {
			return err
		}

		token, err := authToken.Package()
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{"token": token, "type": tokType.String()})
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

		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Minute)
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
			if err == memcache.ErrNotFound {
				return APIError{
					Code:           http.StatusBadRequest,
					PublicMessage:  "Invalid link",
					PrivateMessage: "email token not found in memory cache.",
					InternalError:  err,
				}
			}
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

		isVerified, err := queries.IsUserEmailVerified(ctx, userId)
		if err != nil {
			return err
		}

		if isVerified {
			return c.String(http.StatusOK, "verified")
		}

		err = queries.SetUserEmailVerifiedStatus(
			ctx,
			db.SetUserEmailVerifiedStatusParams{
				ID:            userId,
				EmailVerified: true,
				UpdatedAt:     time.Now().UTC().Format(time.RFC3339Nano),
			},
		)
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

func ResendVerificationEmail(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}

		if user.EmailVerified {
			return c.String(http.StatusBadRequest, "Email already verified.")
		}

		logger := middlewares.GetLogger(c)
		go sendVerificationEmail(user.ID, user.Email, logger)

		return nil
	}
}

// SetupTOTP start TOTP setup for a registered user.
func SetupTOTP(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}

		if user.TotpLocked {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "To re-setup TOTP please remove the TOTP first.",
			}
		}

		issuer := "Konbini"
		accountName := user.Email

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      issuer,
			AccountName: accountName,
		})
		if err != nil {
			return APIError{
				Code:           http.StatusInternalServerError,
				PublicMessage:  "Failed to setup TOTP. Please try again.",
				PrivateMessage: "Error generating TOTP",
				InternalError:  err,
			}
		}

		conn, err := connector.Connect()
		if err != nil {
			return APIError{
				Code:           http.StatusInternalServerError,
				PublicMessage:  "Failed to setup TOTP. Please try again.",
				PrivateMessage: "Error connecting to database",
				InternalError:  err,
			}
		}
		defer conn.Close()

		q := db.New(conn)

		err = q.SetUserTOTPSecret(c.Request().Context(), db.SetUserTOTPSecretParams{
			TotpSecret: sql.NullString{
				String: key.Secret(),
				Valid:  true,
			},
		})
		if err != nil {
			return APIError{
				Code:           http.StatusInternalServerError,
				PublicMessage:  "Failed to setup TOTP. Please try again.",
				PrivateMessage: "Error saving user TOTP secret.",
				InternalError:  err,
			}
		}

		url := key.URL()

		return c.JSON(http.StatusOK, map[string]string{"url": url})
	}
}

type SetupTOTPLockRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// SetupTOTPLock finishes the TOTP setup and generates backup codes for the client.
func SetupTOTPLock(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}

		if user.TotpLocked {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "To re-setup TOTP please remove it first.",
			}
		}

		body, err := middlewares.GetJsonBody[SetupTOTPLockRequest](c)
		if err != nil {
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Failed to get the validated json body",
				InternalError:  err,
			}
		}

		if !totp.Validate(body.Code, user.TotpSecret.String) {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Invalid code.",
			}
		}

		codes := make([]string, 6)
		for i := 0; i < 6; i++ {
			code, err := utils.RandomBytes(16)
			if err != nil {
				return APIError{
					Code:           http.StatusInternalServerError,
					PublicMessage:  "Failed to validate TOTP code.",
					PrivateMessage: "Failed when generating new random bytes for recovery code",
					InternalError:  err,
				}
			}
			codes[i] = hex.EncodeToString(code)
		}

		conn, err := connector.Connect()
		if err != nil {
			return APIError{
				Code:          http.StatusInternalServerError,
				InternalError: err,
			}
		}
		defer conn.Close()

		q := db.New(conn)
		err = q.NewRecoveryCodes(c.Request().Context(), db.NewRecoveryCodesParams{
			UserID:    user.ID,
			CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
			Code:      codes[0],
			Code_2:    codes[1],
			Code_3:    codes[2],
			Code_4:    codes[3],
			Code_5:    codes[4],
			Code_6:    codes[5],
		})
		if err != nil {
			return APIError{
				Code:           http.StatusInternalServerError,
				PublicMessage:  "Failed to verify TOTP code.",
				PrivateMessage: "Failed to store the recovery codes in the database.",
				InternalError:  err,
			}
		}

		return c.JSON(http.StatusOK, map[string][]string{
			"recover_codes": codes,
		})
	}
}

// RemoveTOTP removes the TOTP that has been setup for the requesting user.
func RemoveTOTP(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}

		if !user.TotpLocked || !user.TotpSecret.Valid {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "No TOTP setup.",
			}
		}

		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()
		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		q := db.New(conn)
		q = q.WithTx(tx)

		err = q.RemoveUserRecoveryCodes(c.Request().Context(), user.ID)
		if err != nil {
			tx.Rollback()
			return APIError{
				Code:           http.StatusInternalServerError,
				PublicMessage:  "Failed to remove TOTP",
				PrivateMessage: "Failed to remove the user recovery codes from database",
				InternalError:  err,
			}
		}

		err = q.RemoveUserTOTPSecret(
			c.Request().Context(),
			db.RemoveUserTOTPSecretParams{ID: user.ID, UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano)},
		)
		if err != nil {
			tx.Rollback()
			return APIError{
				Code:           http.StatusInternalServerError,
				PublicMessage:  "Failed to remove TOTP",
				PrivateMessage: "Failed to update the user totp secret and locked properites in database",
				InternalError:  err,
			}
		}

		// invalidate all tokens that has been served to the user
		err = q.DeleteUserJwts(c.Request().Context(), user.ID)
		if err != nil {
			tx.Rollback()
			return APIError{
				Code:           http.StatusInternalServerError,
				PublicMessage:  "Failed to remove TOTP",
				PrivateMessage: "Failed to invalidate all tokens in database",
				InternalError:  err,
			}
		}

		err = tx.Commit()
		if err != nil {
			return APIError{
				Code:           http.StatusInternalServerError,
				PublicMessage:  "Failed to remove TOTP",
				PrivateMessage: "Failed to commit changes to completely remove TOTP",
				InternalError:  err,
			}
		}

		return c.NoContent(http.StatusOK)
	}
}
