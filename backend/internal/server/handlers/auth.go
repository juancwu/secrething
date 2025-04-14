package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/juancwu/secrething/internal/server/api"
	"github.com/juancwu/secrething/internal/server/cookie"
	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/middleware"
	"github.com/juancwu/secrething/internal/server/services"
	"github.com/juancwu/secrething/internal/server/templates"
)

// Error codes for authentication
const (
	// Registration errors
	UserAlreadyExistsCode string = "ERR_USR_EMAIL_DUP_4023"
	RegistrationErrorCode string = "ERR_REGISTER_4024"

	// Login errors
	InvalidCredentialsCode     string = "ERR_AUTH_INVALID_CREDS_4010"
	UserNotFoundCode           string = "ERR_AUTH_USER_NOT_FOUND_4011"
	UserAccountLockedCode      string = "ERR_AUTH_ACCOUNT_LOCKED_4012"
	UserEmailNotVerifiedCode   string = "ERR_AUTH_EMAIL_NOT_VERIFIED_4013"
	AuthenticationRequiredCode string = "ERR_AUTH_REQUIRED_4014"

	// TOTP errors
	RequiresTotpCode     string = "ERR_AUTH_TOTP_REQUIRED_4020"
	InvalidTOTPCodeCode  string = "ERR_AUTH_INVALID_TOTP_4021"
	InvalidTOTPTokenCode string = "ERR_AUTH_INVALID_TOTP_TOKEN_4022"

	// Token errors
	InvalidRefreshTokenCode     string = "ERR_AUTH_INVALID_REFRESH_4030"
	RefreshTokenRequiredCode    string = "ERR_AUTH_REFRESH_REQUIRED_4031"
	UnauthorizedTokenAccessCode string = "ERR_AUTH_UNAUTHORIZED_TOKEN_4032"
)

// Default timeout for auth endpoints
const defaultTimeout = time.Minute

// Cookie settings for the refresh token
const (
	RefreshTokenCookieName = "refresh_token"
	RefreshTokenCookiePath = "/api/auth"
	CookieMaxAge           = 7 * 24 * 60 * 60 // 7 days in seconds
)

type AuthHandler struct{}

func NewAuthHandler() AuthHandler {
	return AuthHandler{}
}

// ConfigureRoutes implements the Handler interface
func (h AuthHandler) ConfigureRoutes(e *echo.Echo, v *validator.Validator) {
	e.POST("/api/auth/sign-up", h.createUser, middleware.SetValidator(v, getCreateUserRequestMessages()))
	e.POST("/api/auth/sign-in", h.signIn, middleware.SetValidator(v, getSignInRequestMessages()))
	e.POST("/api/auth/refresh", h.refreshToken)

	e.POST("/api/auth/cli/sign-up", h.createUser, middleware.SetValidator(v, getCreateUserRequestMessages()))
	e.POST("/api/auth/cli/sign-in", h.signIn, middleware.SetValidator(v, getSignInRequestMessages()))
	e.POST("/api/auth/cli/refresh", h.refreshToken, middleware.SetValidator(v, getRefreshRequestMessages()))

	// The account activate route verifies the token in the query string hence it is not protected
	// Get method so that when the user clicks on the verify button or paste the url in the browser,
	// verification can happen immediately.
	e.GET("/api/auth/account/activate", h.verifyEmail)

	// Protected routes
	protected := e.Group("", middleware.Protected())
	protected.POST("/api/auth/account/resend-activation", h.resendAccountActivation)
}

type createUserRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,password"`
	Name     *string `json:"name" validate:"omitnil,omitempty,max=50"`
}

type userResponse struct {
	UserID        db.UserID `json:"uid"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Name          *string   `json:"name"`
}

type createUserResponse struct {
	AccessToken  string       `json:"atk"`
	RefreshToken string       `json:"rtk"`
	User         userResponse `json:"user"`
}

func getCreateUserRequestMessages() validator.ValidationMessages {
	vm := validator.NewValidationMessages()
	vm.SetMessage("email", "required", "Email is required.")
	vm.SetMessage("email", "email", "Expecting a valid email, but received '{1}'.")
	vm.SetMessage("password", "required", "Password is required.")
	vm.SetMessage("name", "max", "Name must be at most {2} characters long.")
	return vm
}

// createUser handles requests to create new users.
func (AuthHandler) createUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""

	var body createUserRequest
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := middleware.Validate(c, &body); err != nil {
		return err
	}

	authService := services.NewAuthService()

	user, err := authService.CreateUser(ctx, body.Email, body.Password, body.Name)
	if err != nil {
		if serviceErr, ok := err.(services.AuthServiceError); ok && serviceErr.IsType(services.AuthServiceErrUserAlreadyExists) {
			return api.NewBadRequestError(serviceErr.Error(), UserAlreadyExistsCode, requestID, err)
		}
		return api.NewInternalServerError("Failed to create user", requestID, err)
	}

	// Determine client type based on the endpoint path
	clientType := services.ClientTypeWeb
	if c.Path() == "/api/auth/cli/sign-up" {
		clientType = services.ClientTypeCLI
	}

	tokenService := services.NewTokenService()
	tokenPair, err := tokenService.GenerateTokenPair(ctx, user.UserID, clientType)
	if err != nil {
		// TODO: logging
		fmt.Println(err.Error())
		return c.JSON(201, createUserResponse{
			User: userResponse{
				UserID:        user.UserID,
				Email:         user.Email,
				EmailVerified: user.EmailVerified,
				Name:          user.Name,
			},
		})
	}

	go func(email string, userID db.UserID) {
		ctx, cancel := context.WithTimeoutCause(context.Background(), time.Second*30, fmt.Errorf("Send account verification email to '%s' timeout", email))
		defer cancel()

		token, err := tokenService.NewAccountActivateToken(ctx, userID)
		if err != nil {
			fmt.Printf("Failed to send account verification email: %v\n", err)
			return
		}

		emailService := services.NewEmailService()
		emailService.SendAccountVerificationEmail(ctx, email, token)
	}(user.Email, user.UserID)

	resBody := createUserResponse{
		AccessToken: tokenPair.AccessToken,
		User: userResponse{
			UserID:        user.UserID,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Name:          user.Name,
		},
	}

	if clientType == services.ClientTypeWeb {
		cookie.SetRefreshToken(c, tokenPair.RefreshToken)
	} else {
		resBody.RefreshToken = tokenPair.RefreshToken
	}

	return c.JSON(http.StatusCreated, resBody)
}

type signInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func getSignInRequestMessages() validator.ValidationMessages {
	vm := validator.NewValidationMessages()
	vm.SetMessage("email", "required", "Email is required.")
	vm.SetMessage("email", "email", "Expecting a valid email, but received '{1}'.")
	vm.SetMessage("password", "required", "Password is required.")
	return vm
}

// signIn handles user authentication and returns tokens
func (AuthHandler) signIn(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""

	var body signInRequest
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := middleware.Validate(c, &body); err != nil {
		return err
	}

	authService := services.NewAuthService()
	user, err := authService.AuthenticateUser(ctx, body.Email, body.Password)
	if err != nil {
		if serviceErr, ok := err.(services.AuthServiceError); ok {
			switch serviceErr.Type {
			case services.AuthServiceErrUserNotFound:
				return api.NewNotFoundError(serviceErr.Error(), UserNotFoundCode, requestID, serviceErr)
			case services.AuthServiceErrInvalidCredentials:
				return api.NewUnauthorizedError(serviceErr.Error(), InvalidCredentialsCode, requestID, serviceErr)
			case services.AuthServiceErrAccountLocked:
				return api.NewForbiddenError(serviceErr.Error(), UserAccountLockedCode, requestID, serviceErr)
			case services.AuthServiceErrEmailNotVerified:
				return api.NewForbiddenError(serviceErr.Error(), UserEmailNotVerifiedCode, requestID, serviceErr)
			}
		}
		return api.NewInternalServerError("Failed to authenticate user", requestID, err)
	}

	// Determine client type based on the endpoint path
	clientType := services.ClientTypeWeb
	if c.Path() == "/api/auth/cli/sign-in" {
		clientType = services.ClientTypeCLI
	}

	tokenService := services.NewTokenService()
	tokenPair, err := tokenService.GenerateTokenPair(ctx, user.UserID, clientType)
	if err != nil {
		return api.NewInternalServerError("Failed to generate token pair", requestID, err)
	}

	resBody := createUserResponse{
		AccessToken: tokenPair.AccessToken,
		User: userResponse{
			UserID:        user.UserID,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Name:          user.Name,
		},
	}

	if clientType == services.ClientTypeWeb {
		cookie.SetRefreshToken(c, tokenPair.RefreshToken)
	} else {
		resBody.RefreshToken = tokenPair.RefreshToken
	}

	return c.JSON(http.StatusOK, resBody)
}

type refreshRequest struct {
	Value string `json:"refresh_token" validate:"required"`
}

func getRefreshRequestMessages() validator.ValidationMessages {
	vm := validator.NewValidationMessages()
	vm.SetMessage("refresh_token", "required", "Refresh token is required")
	return vm
}

// refreshToken handles token refresh using a refresh token
func (AuthHandler) refreshToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	requestID := ""

	var refreshToken string
	if c.Path() == "/api/auth/cli/refresh" {
		// Get refresh token from request body
		var reqBody refreshRequest
		if err := c.Bind(&reqBody); err != nil {
			return err
		}
		if err := middleware.Validate(c, &reqBody); err != nil {
			return err
		}
		refreshToken = reqBody.Value
	} else {
		// Get refresh token from cookie
		cookieObj, err := c.Cookie(cookie.RefreshTokenKey)
		if err != nil {
			return api.NewBadRequestError("Refresh token cookie not found", RefreshTokenRequiredCode, requestID, err)
		}

		refreshToken = cookieObj.Value
	}

	if refreshToken == "" {
		return api.NewBadRequestError("Refresh token is empty", RefreshTokenRequiredCode, requestID, nil)
	}

	// Verify refresh token
	tokenService := services.NewTokenService()
	payload, err := tokenService.VerifyToken(ctx, refreshToken, services.StdPackage)
	if err != nil {
		switch err := err.(type) {
		case services.TokenServiceError:
			switch err.Type {
			case services.TokenServiceErrExpired:
				// Clear the invalid refresh token cookie
				cookie.UnsetRefreshToken(c)
				return api.NewUnauthorizedError("Refresh token has expired, please sign in again", InvalidRefreshTokenCode, requestID, err)
			case services.TokenServiceErrInvalid, services.TokenServiceErrDecryption:
				return api.NewBadRequestError(err.Error(), InvalidRefreshTokenCode, requestID, err)
			}
		}
		return api.NewInternalServerError("Failed to validate refresh token", requestID, err)
	}

	if payload.TokenType != services.TokenTypeRefresh {
		return api.NewBadRequestError("Invalid token type, expected refresh token", InvalidRefreshTokenCode, requestID, nil)
	}

	// Verify token in database
	q, err := db.Query()
	if err != nil {
		return api.NewInternalServerError("Database error", requestID, err)
	}

	// Get user from database
	user, err := q.GetUserByID(ctx, payload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return api.NewUnauthorizedError("User not found", UserNotFoundCode, requestID, err)
		}
		return api.NewInternalServerError("Failed to get user", requestID, err)
	}

	// Generate new token pair
	tokenPair, err := tokenService.GenerateTokenPair(ctx, user.UserID, payload.ClientType)
	if err != nil {
		return api.NewInternalServerError("Failed to generate new token pair", requestID, err)
	}

	// Set new refresh token cookie for web clients
	if payload.ClientType == services.ClientTypeWeb {
		cookie.SetRefreshToken(c, tokenPair.RefreshToken)
		// Only return access token in response for web clients
		return c.JSON(http.StatusOK, map[string]interface{}{
			"access_token": tokenPair.AccessToken,
			"user": userResponse{
				UserID:        user.UserID,
				Email:         user.Email,
				EmailVerified: user.EmailVerified,
				Name:          user.Name,
			},
		})
	}

	// Return both tokens for CLI clients
	return c.JSON(http.StatusOK, map[string]interface{}{
		"tokens": tokenPair,
		"user": userResponse{
			UserID:        user.UserID,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Name:          user.Name,
		},
	})
}

func (h *AuthHandler) resendAccountActivation(c echo.Context) error {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeoutCause(context.Background(), time.Second*30, fmt.Errorf("Send account verification email to '%s' timeout", user.Email))
	defer cancel()

	q, err := db.Query()
	if err != nil {
		return err
	}

	if err := q.DeleteTokensByType(ctx, db.DeleteTokensByTypeParams{UserID: user.UserID, TokenType: services.TokenTypeAccountActivate}); err != nil {
		return err
	}

	tokenService := services.NewTokenService()
	token, err := tokenService.NewAccountActivateToken(ctx, user.UserID)
	if err != nil {
		return err
	}

	emailService := services.NewEmailService()
	if err := emailService.SendAccountVerificationEmail(ctx, user.Email, token); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *AuthHandler) verifyEmail(c echo.Context) error {
	// Get token from query parameter and make sure that it is defined
	token := c.QueryParam("token")
	// If token is empty, then render general something went wrong page
	if token == "" {
		return c.HTML(http.StatusOK, templates.StaticInternalError)
	}

	tokenService := services.NewTokenService()
	authService := services.NewAuthService()

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// Verify the token
	payload, err := tokenService.VerifyToken(ctx, token, services.UrlPackage)
	if err != nil {
		fmt.Println(err)
		return c.HTML(http.StatusOK, templates.StaticInternalError)
	}

	// If token is not of expected type, then fail verification
	if payload.TokenType != services.TokenTypeAccountActivate {
		return c.HTML(http.StatusOK, templates.StaticVerificationError)
	}

	exists, err := authService.ExistsUserWithID(ctx, payload.UserID)
	if err != nil {
		// TODO: log error in sentry
		fmt.Println(err)
		return c.HTML(http.StatusOK, templates.StaticInternalError)
	}
	if !exists {
		return c.HTML(http.StatusOK, templates.StaticVerificationError)
	}

	// Check that the token exists in the database
	_, err = tokenService.GetTokenByID(ctx, payload.TokenID)
	if err != nil {
		if db.IsNoRows(err) {
			return c.HTML(http.StatusOK, templates.StaticVerificationError)
		}
		fmt.Println(err)
		return c.HTML(http.StatusOK, templates.StaticInternalError)
	}

	// Set the email verification status to true
	_, err = authService.UpdateUserEmailVerification(ctx, payload.UserID, true)
	if err != nil {
		// TODO: log error in sentry
		fmt.Println(err)
		return c.HTML(http.StatusOK, templates.StaticVerificationError)
	}

	return c.HTML(http.StatusOK, templates.StaticVerificationSuccess)
}
