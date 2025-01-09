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
	"github.com/rs/zerolog/log"
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

		now := time.Now()
		userId, err := queries.CreateUser(ctx, db.CreateUserParams{
			Email:     body.Email,
			Password:  hash,
			Nickname:  body.NickName,
			CreatedAt: utils.FormatRFC3339NanoFixed(now),
			UpdatedAt: utils.FormatRFC3339NanoFixed(now),
		})
		if err != nil {
			return err
		}

		logger.Info().Str("user_id", userId).Msg("New user registered.")
		go sendVerificationEmail(userId, body.Email, logger)

		// generate a partial token so that the user can immediately setup TOTP
		exp := now.Add(time.Hour * 24 * 7)
		var authToken *services.AuthToken
		dbJwt, err := queries.NewAuthToken(ctx, db.NewAuthTokenParams{
			UserID:    userId,
			TokenType: services.PARTIAL_USER_TOKEN_TYPE.String(),
			CreatedAt: utils.FormatRFC3339NanoFixed(now),
			ExpiresAt: utils.FormatRFC3339NanoFixed(exp),
		})
		if err != nil {
			return err
		}
		authToken, err = services.NewAuthToken(dbJwt.ID, userId, services.PARTIAL_USER_TOKEN_TYPE, exp)
		if err != nil {
			return err
		}

		token, err := authToken.Package()
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]string{"token": token, "type": services.PARTIAL_USER_TOKEN_TYPE.String()})
	}
}

type LoginRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required"`
	TOTPCode *string `json:"totp_code,omitempty" validate:"omitnil,omitempty,required,len=6|len=32"`
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
		if user.TotpSecret == nil || !user.EmailVerified {
			tokType = services.PARTIAL_USER_TOKEN_TYPE
		} else if user.TotpLocked && user.TotpSecret != nil {
			if body.TOTPCode == nil {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "User has TOTP setup, code is required. Make a new login request with the totp_code field in the body.",
				}
			}

			// recovery code
			if len(*body.TOTPCode) == 32 {
				// check database for recovery code
				recoveryCode, err := queries.GetRecoveryCode(ctx, db.GetRecoveryCodeParams{
					UserID: user.ID,
					Code:   *body.TOTPCode,
				})
				if err != nil {
					if err == sql.ErrNoRows {
						return APIError{
							Code:           http.StatusBadRequest,
							PublicMessage:  "Invalid recovery code.",
							PrivateMessage: "No recovery code found in database.",
							InternalError:  err,
						}
					}
					return APIError{
						Code:           http.StatusInternalServerError,
						PrivateMessage: "Failed to get recovery code from database.",
						InternalError:  err,
					}
				}

				if recoveryCode.Used {
					return APIError{
						Code:           http.StatusBadRequest,
						PublicMessage:  "Recovery code has been used before. Please use another one.",
						PrivateMessage: "Used recovery code. Reject.",
					}
				}

				err = queries.UseRecoveryCode(ctx, db.UseRecoveryCodeParams{
					UserID: user.ID,
					Code:   recoveryCode.Code,
				})
				if err != nil {
					return APIError{
						Code:           http.StatusInternalServerError,
						PrivateMessage: "Failed to use recovery code.",
						InternalError:  err,
					}
				}
			} else if !totp.Validate(*body.TOTPCode, *user.TotpSecret) {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "Invalid TOTP code.",
				}
			}

			tokType = services.FULL_USER_TOKEN_TYPE
		} else {
			tokType = services.FULL_USER_TOKEN_TYPE
		}

		authToken, err := newAuthToken(ctx, queries, user.ID, tokType)
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
				UpdatedAt:     utils.FormatRFC3339NanoFixed(time.Now()),
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

		secret := key.Secret()
		err = q.SetUserTOTPSecret(c.Request().Context(), db.SetUserTOTPSecretParams{
			TotpSecret: &secret,
			UpdatedAt:  utils.FormatRFC3339NanoFixed(time.Now()),
			ID:         user.ID,
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
		logger := middlewares.GetLogger(c)
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}

		if user.TotpSecret == nil {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "No TOTP setup yet.",
			}
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

		log.Debug().Str("code", body.Code).Send()

		if !totp.Validate(body.Code, *user.TotpSecret) {
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

		tx, err := conn.Begin()
		if err != nil {
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Failed to start transtaction",
				InternalError:  err,
			}
		}

		q := db.New(conn)
		q = q.WithTx(tx)

		err = q.NewRecoveryCodes(c.Request().Context(), db.NewRecoveryCodesParams{
			UserID:    user.ID,
			CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
			Code:      codes[0],
			Code_2:    codes[1],
			Code_3:    codes[2],
			Code_4:    codes[3],
			Code_5:    codes[4],
			Code_6:    codes[5],
		})
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error().Err(err).Msg("Failed to rollback")
			}
			return APIError{
				Code:           http.StatusInternalServerError,
				PublicMessage:  "Failed to verify TOTP code.",
				PrivateMessage: "Failed to store the recovery codes in the database.",
				InternalError:  err,
			}
		}

		err = q.LockUserTOTP(c.Request().Context(), db.LockUserTOTPParams{
			UpdatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
			ID:        user.ID,
		})
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error().Err(err).Msg("Failed to rollback")
			}
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Failed to lock user TOTP status in database",
				InternalError:  err,
			}
		}

		var token string
		if user.EmailVerified {
			// generate a full token for the user to start using instead of the partial token
			authToken, err := newAuthToken(c.Request().Context(), q, user.ID, services.FULL_USER_TOKEN_TYPE)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					logger.Error().Err(err).Msg("Failed to rollback")
				}
				return APIError{
					Code:           http.StatusInternalServerError,
					PrivateMessage: "Failed to create new auth token",
					InternalError:  err,
				}
			}

			token, err = authToken.Package()
			if err != nil {
				if err := tx.Rollback(); err != nil {
					logger.Error().Err(err).Msg("Failed to rollback")
				}
				return APIError{
					Code:           http.StatusInternalServerError,
					PrivateMessage: "Failed to package auth token",
					InternalError:  err,
				}
			}

			// remove all partial tokens, make them invalid
			err = q.DeleteAllTokensByTypeAndUserID(c.Request().Context(), db.DeleteAllTokensByTypeAndUserIDParams{
				UserID:    user.ID,
				TokenType: services.PARTIAL_USER_TOKEN_TYPE.String(),
			})
			if err != nil {
				if err := tx.Rollback(); err != nil {
					logger.Error().Err(err).Msg("Failed to rollback")
				}
				return APIError{
					Code:           http.StatusInternalServerError,
					PrivateMessage: "Failed to delete all partial tokens owned by user",
					InternalError:  err,
				}
			}
		}

		err = tx.Commit()
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error().Err(err).Msg("Failed to rollback")
			}
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Failed to commit transaction",
				InternalError:  err,
			}
		}

		if token != "" {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"recover_codes": codes,
				"token":         token,
				"type":          services.FULL_USER_TOKEN_TYPE.String(),
			})
		}

		return c.JSON(http.StatusOK, map[string][]string{
			"recover_codes": codes,
		})
	}
}

// RemoveTOTPRequest is the expected request body for remove totp request
type RemoveTOTPRequest struct {
	Code string `json:"code" validate:"required,len=6|len=32|len=8"`
}

// RemoveTOTP removes the TOTP that has been setup for the requesting user.
func RemoveTOTP(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}

		body, err := middlewares.GetJsonBody[RemoveTOTPRequest](c)
		if err != nil {
			return err
		}

		if !user.TotpLocked || user.TotpSecret == nil {
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

		switch len(body.Code) {
		case 6:
			if !totp.Validate(body.Code, *user.TotpSecret) {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "Invalid TOTP code.",
				}
			}
		case 8:
			// this is a code sent through email
			if !user.EmailVerified {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "Email must be verified to use email codes.",
				}
			}
			k, exp, found := memcache.Cache().GetWithExpiration("email_code_" + user.ID)
			if !found || time.Now().After(exp) {
				return APIError{
					Code:           http.StatusBadRequest,
					PublicMessage:  "Invalid 2FA code.",
					PrivateMessage: "Email code not found in memory cache or expired",
				}
			}
			emailCode, ok := k.(string)
			if !ok {
				return APIError{
					Code:           http.StatusInternalServerError,
					PrivateMessage: "Failed to cast email code stored in memory cache to string",
				}
			}
			if body.Code != emailCode {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "Invalid 2FA code.",
				}
			}
		case 32:
			q := db.New(conn)
			if err := verifyRecoveryCode(c.Request().Context(), q, user.ID, body.Code); err != nil {
				return APIError{
					Code:           http.StatusBadRequest,
					PublicMessage:  "Invalid TOTP code.",
					PrivateMessage: "Verification failed at recovery code.",
					InternalError:  err,
				}
			}
		default:
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Invalid 'code' in body.",
			}
		}

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
			db.RemoveUserTOTPSecretParams{
				ID:        user.ID,
				UpdatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
			},
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
		err = q.DeleteUserAuthTokens(c.Request().Context(), user.ID)
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
