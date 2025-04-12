package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/juancwu/secrething/internal/server/api"
	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/middleware"
	"github.com/juancwu/secrething/internal/server/services"
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
func (h AuthHandler) ConfigureRoutes(g *echo.Group, v *validator.Validator) {
	g.POST("/sign-up", h.createUser, middleware.SetValidator(v, getCreateUserRequestMessages()))
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
	TokenPair services.TokenPair `json:"tokens"`
	User      userResponse       `json:"user"`
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
		if serviceErr, ok := err.(services.AuthServiceError); ok && serviceErr.IsType(services.ErrUserAlreadyExists) {
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
		return c.JSON(201, createUserResponse{
			User: userResponse{
				UserID:        user.UserID,
				Email:         user.Email,
				EmailVerified: user.EmailVerified,
				Name:          user.Name,
			},
		})
	}

	go func(clientType string, email string, userID db.UserID) {
		ctx, cancel := context.WithTimeoutCause(context.Background(), time.Second*30, fmt.Errorf("Send account verification email to '%s' timeout", email))
		defer cancel()

		now := time.Now()
		exp := now.Add(time.Hour * 24)

		emailToken, err := tokenService.GenerateToken(ctx, userID, services.TokenAccountActivate, clientType, now, exp)
		if err != nil {
			fmt.Printf("Failed to send account verification email: %v\n", err)
			return
		}

		emailService := services.NewEmailService()
		emailService.SendAccountVerificationEmail(ctx, email, emailToken.TokenID)
	}(clientType, user.Email, user.UserID)

	return c.JSON(201, createUserResponse{
		TokenPair: tokenPair,
		User: userResponse{
			UserID:        user.UserID,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Name:          user.Name,
		},
	})
}
