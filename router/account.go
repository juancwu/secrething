package router

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/juancwu/konbini/store"
	"github.com/juancwu/konbini/utils"
	"github.com/juancwu/konbini/views"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

const (
	HeaderXRefreshToken string = "X-Refresh-Token"
)

// SetupAccountRoutes setups the account related routes. These routes belong to /api/v1
func SetupAccountRoutes(e RouteGroup) {
	e.POST("/account/signup", handleSignup, useValidateRequestBody(signupRequest{}))
	e.POST("/account/login", handleLogin, useValidateRequestBody(loginRequest{}))
	e.GET("/account/new-token", handleNewToken)
	e.GET("/account/send-verification-email", handleSendVerificationEmail)
	e.GET("/account/verify-email", handleVerifyEmail)
}

func handleSignup(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	body, ok := c.Get("body").(*signupRequest)
	if !ok {
		logger.Error("Invalid body type", zap.String("request_id", requestId))
		return c.JSON(
			http.StatusInternalServerError,
			apiResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    http.StatusText(http.StatusInternalServerError),
			},
		)
	}

	exists, err := store.UserExists(body.Email)
	if err != nil {
		logger.Error(
			"Failed to get user by email",
			zap.String("request_id", requestId),
			zap.Error(err),
		)
		return c.JSON(
			http.StatusInternalServerError,
			apiResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    http.StatusText(http.StatusInternalServerError),
			},
		)
	}
	if exists {
		logger.Error(
			"User already exists",
			zap.String("email", body.Email),
			zap.String("request_id", requestId),
		)
		return c.JSON(
			http.StatusBadRequest,
			apiResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "User with given email already exists.",
			},
		)
	}

	// create new user
	userId, err := store.CreateUser(body.Email, body.Password, body.FirstName, body.LastName)
	if err != nil {
		logger.Error(
			"Failed to create new user.",
			zap.String("email", body.Email),
			zap.String("request_id", requestId),
			zap.Error(err),
		)
		return c.JSON(
			http.StatusInternalServerError,
			apiResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    http.StatusText(http.StatusInternalServerError),
			},
		)
	}

	// send verification email using a go routing to not block the response
	go sendVerificationEmail(body.Email, body.FirstName, userId)

	return c.JSON(http.StatusCreated, apiResponse{
		StatusCode: http.StatusCreated,
		Message:    "We have sent an email to you. Please verify your email.",
	})
}

func handleLogin(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	body, ok := c.Get("body").(*loginRequest)
	if !ok {
		logger.Error("Invalid body type", zap.String("request_id", requestId))
		return writeApiErrorJSON(c, requestId)
	}

	user, err := store.GetUserWithPasswordValidation(body.Email, body.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("No user found with given email and password.", zap.String("email", body.Email), zap.String("request_id", requestId))
			return c.JSON(
				http.StatusBadRequest,
				apiResponse{
					StatusCode: http.StatusBadRequest,
					Message:    fmt.Sprintf("invalid credentials (%s)", requestId),
				},
			)
		}
		logger.Error("Failed to get user with email and password.", zap.Error(err), zap.String("request_id", requestId))
		return writeApiErrorJSON(c, requestId)
	}

	if !user.EmailVerified {
		logger.Error("User login attempt when email has not been verified.", zap.String("email", user.Email), zap.String("request_id", requestId))
		return c.JSON(
			http.StatusBadRequest,
			apiResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Please verify your email before logging in.",
			},
		)
	}

	// get signed jwt to send back to user
	// should generate two tokens, refresh and access token
	accessToken, err := generateAccessToken(user.Id)
	if err != nil {
		logger.Error("Failed to generate access token", zap.Error(err), zap.String("request_id", requestId))
		return writeApiErrorJSON(c, requestId)
	}
	refreshToken, err := generateRefreshToken(user.Id)
	if err != nil {
		logger.Error("Failed to generate refresh token", zap.Error(err), zap.String("request_id", requestId))
		return writeApiErrorJSON(c, requestId)
	}

	return c.JSON(
		http.StatusOK,
		loginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	)
}

// handleNewToken handles when generating a new access token using a refresh token.
func handleNewToken(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	refrestTokenString := c.Request().Header.Get(HeaderXRefreshToken)
	if refrestTokenString == "" {
		return c.JSON(
			http.StatusBadRequest,
			apiResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Missing required header 'X-Refresh-Token'",
			},
		)
	}

	// verify the token
	token, err := verifyJWT(refrestTokenString)
	if err != nil {
		if err.Error() == fmt.Sprintf("%s: %s", jwt.ErrTokenInvalidClaims.Error(), jwt.ErrTokenExpired.Error()) {
			return c.JSON(
				http.StatusBadRequest,
				apiResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Refresh token expired. Login again to get a new one.",
				},
			)
		}
		logger.Error("Failed to verify refresh token.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
		return writeApiErrorJSON(c, requestId)

	}

	// generate a new access token
	claims := token.Claims.(*jwtAuthClaims)
	accessToken, err := generateAccessToken(claims.ID)
	if err != nil {
		logger.Error("Failed to generate a new access token.", zap.Error(err), zap.String("request_id", requestId))
		return writeApiErrorJSON(c, requestId)
	}

	return c.JSON(
		http.StatusCreated,
		newTokenResponse{
			AccessToken: accessToken,
		},
	)
}

func handleSendVerificationEmail(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)

	claims, err := useJWT(c, JWT_ACCESS_TOKEN_TYPE)
	if err != nil {
		logger.Error("Failed to validate jwt.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
		return writeUnauthorized(c, requestId)
	}

	user, err := store.GetUserWithId(claims.UserId)
	if err != nil {
		logger.Error("Failed to get user id with email", zap.String("user_id", claims.UserId), zap.String("request_id", requestId), zap.Error(err))
		return writeApiErrorJSON(c, requestId)
	}

	if user.EmailVerified {
		return c.JSON(
			http.StatusOK,
			apiResponse{
				StatusCode: http.StatusOK,
				Message:    "User email has already been verified.",
			},
		)
	}

	_, err = store.DeleteAllEmailVerificationFromUser(user.Id)
	if err != nil {
		logger.Error("Failed to delete all existing email verifications linked to requesting user.")
	}

	// send new email
	go sendVerificationEmail(user.Email, user.FirstName, user.Id)

	return c.JSON(
		http.StatusOK,
		apiResponse{
			StatusCode: http.StatusOK,
			Message:    "Verification email scheduled to send. You should receive an email from us shortly.",
		},
	)
}

func handleVerifyEmail(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(
			http.StatusBadRequest,
			apiResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Missing required query parameter 'code'.",
			},
		)
	}
	now := time.Now()

	// get email verification record in database with code
	emailVerification, err := store.GetEmailVerificationWithCode(code)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(
				http.StatusNotFound,
				apiResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Code does not exists. Please get a new code.",
				},
			)
		}
		logger.Error("Failed to get email verification record.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
		return writeApiErrorJSON(c, requestId)
	}

	if now.After(emailVerification.ExpiresAt) {
		go func() {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			logger.Info("Updating expired email verification status.", zap.Int64("email_verification_id", emailVerification.Id), zap.String(echo.HeaderXRequestID, requestId))
			err := store.DeleteEmailVerification(emailVerification.Id)
			if err != nil {
				logger.Error("Failed to update expired email verification status.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
				return
			}
		}()
		return c.JSON(
			http.StatusBadRequest,
			apiResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Code has expired. Please get a new email verification with a new code.",
			},
		)
	}

	tx, err := store.StartTx()
	if err != nil {
		logger.Error("Failed to start transaction to update email verification status.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
		return writeApiErrorJSON(c, requestId)
	}

	// delete the record because there is no need to keep it anymore
	err = store.DeleteEmailVerificationTx(tx, emailVerification.Id)
	if err != nil {
		logger.Error("Failed to delete email verification.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
		go func() {
			err := tx.Rollback()
			if err != nil {
				logger.Error("Failed to rollback transaction.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
			}
		}()
		return writeApiErrorJSON(c, requestId)
	}

	err = store.SetUserEmailVerifiedStatus(tx, emailVerification.UserId, true)
	if err != nil {
		logger.Error("Failed to update user email verified status.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
		go func(tx *sql.Tx, requestId string) {
			err := tx.Rollback()
			if err != nil {
				logger.Error("Failed to rollback transaction.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
			}
		}(tx, requestId)
		return writeApiErrorJSON(c, requestId)
	}

	err = tx.Commit()
	if err != nil {
		logger.Error("Failed to commit transaction changes.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
		return writeApiErrorJSON(c, requestId)
	}

	return c.JSON(
		http.StatusOK,
		apiResponse{
			StatusCode: http.StatusOK,
			Message:    "Thanks for verifying your email.",
		},
	)
}

// From here on, all code are just helper functions but not route handlers.

// sendVerificationEmail is a helper function that sends a verification email.
func sendVerificationEmail(email string, firstName string, userId string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// generate code
	code, err := gonanoid.Generate(store.EMAIL_VERIFICATION_CODE_CHR_POOL, store.EMAIL_VERIFICATION_CODE_LEN)
	if err != nil {
		logger.Error("Failed to generate email verification code on new user created.", zap.Error(err))
		return
	}

	// try to send email first
	var html bytes.Buffer
	err = views.VerifyEmail(firstName, fmt.Sprintf("%s/api/v1/account/verify-email?code=%s", os.Getenv("SERVER_URL"), code)).Render(context.Background(), &html)
	if err != nil {
		logger.Error("Failed to render email verification view on new user created.", zap.Error(err))
		return
	}

	// save the email verification in the database
	_, err = store.CreateEmailVerification(code, userId)
	if err != nil {
		logger.Error("Failed to save email verification in database on new user created.", zap.Error(err))
		return
	}

	// send email
	_, err = utils.SendEmail(os.Getenv("NOREPLY_EMAIL"), []string{email}, "[Konbini] Verify Your Email", html.String())
	if err != nil {
		logger.Error("Failed to send email verification on new user created.", zap.Error(err))
		return
	}
}
