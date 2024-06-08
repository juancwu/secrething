package router

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/juancwu/konbini/store"
	"github.com/juancwu/konbini/utils"
	"github.com/juancwu/konbini/views"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

// SetupAccountRoutes setups the account related routes.
func SetupAccountRoutes(e RouteGroup) {
	e.POST("/account/signup", handleSignup, useValidateRequestBody(signupRequest{}))
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
	go func() {
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		// generate code
		code, err := gonanoid.New(store.EMAIL_VERIFICATION_CODE_LEN)
		if err != nil {
			logger.Error("Failed to generate email verification code on new user created.", zap.Error(err))
			return
		}

		// try to send email first
		var html bytes.Buffer
		err = views.VerifyEmail(*body.FirstName, fmt.Sprintf("%s/api/v1/account/verify-email?code=%s", os.Getenv("SERVER_URL"), code)).Render(context.Background(), &html)
		if err != nil {
			logger.Error("Failed to render email verification view on new user created.", zap.Error(err))
			return
		}

		// send email
		emailId, err := utils.SendEmail(os.Getenv("NOREPLY_EMAIL"), []string{body.Email}, "[Konbini] Verify Your Email", html.String())
		if err != nil {
			logger.Error("Failed to send email verification on new user created.", zap.Error(err))
			return
		}

		// save the email verification in the database
		err = store.CreateEmailVerification(code, userId, emailId)
		if err != nil {
			logger.Error("Failed to save email verification in database on new user created.", zap.Error(err))
			return
		}
	}()

	return c.JSON(http.StatusCreated, apiResponse{
		StatusCode: http.StatusCreated,
		Message:    "We have sent an email to you. Please verify your email.",
	})
}
