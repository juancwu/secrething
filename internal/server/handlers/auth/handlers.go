package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/juancwu/secrething/internal/server/db"
	handlerErrors "github.com/juancwu/secrething/internal/server/handlers/errors"
	"github.com/juancwu/secrething/internal/server/middleware"
	authService "github.com/juancwu/secrething/internal/server/services/auth"
	"github.com/juancwu/secrething/internal/server/utils"
)

// Default timeout for auth endpoints
const defaultTimeout = time.Minute

// Cookie settings for the refresh token
const (
	RefreshTokenCookieName = "refresh_token"
	RefreshTokenCookiePath = "/api/auth"
	CookieMaxAge           = 7 * 24 * 60 * 60 // 7 days in seconds
)

// createUser handles requests to create new users.
func createUser(c echo.Context) error {
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

	user, err := authService.CreateUser(ctx, body.Email, body.Password, body.Name)
	if err != nil {
		if serviceErr, ok := err.(authService.AuthServiceError); ok && serviceErr.Is(authService.UserAlreadyExistsErr) {
			return handlerErrors.NewBadRequest(serviceErr.Error(), UserAlreadyExistsCode, requestID, err)
		}
		return handlerErrors.NewInternalServerError("Failed to create user", requestID, err)
	}

	// Determine client type based on the endpoint path
	clientType := authService.ClientTypeWeb
	if c.Path() == "/api/auth/cli/sign-up" {
		clientType = authService.ClientTypeCLI
	}

	// Generate and return tokens automatically after successful registration
	return generateAndReturnTokens(c, ctx, user, clientType, requestID)
}

// loginRequest is an interface that encompasses both signinRequest and totpVerifyRequest
type loginRequest interface {
	GetEmail() string
	GetPassword() string
}

// GetEmail implements loginRequest interface
func (r *signinRequest) GetEmail() string {
	return r.Email
}

// GetPassword implements loginRequest interface
func (r *signinRequest) GetPassword() string {
	return r.Password
}

// loginUser handles requests to sign-in users.
func loginUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""

	// Parse and validate request body
	var body signinRequest
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := middleware.Validate(c, &body); err != nil {
		return err
	}

	// Get the user by email
	q, err := db.Query()
	if err != nil {
		return handlerErrors.NewInternalServerError("Database error", requestID, err)
	}

	user, err := q.GetUserByEmail(ctx, body.Email)
	if err != nil {
		if db.IsNoRows(err) {
			return handlerErrors.NewUnauthorizedError("Invalid email or password", InvalidCredentialsCode, requestID, nil)
		}
		return handlerErrors.NewInternalServerError("Failed to retrieve user", requestID, err)
	}

	// Verify the password
	if match, err := utils.VerifyPassword(body.Password, user.PasswordHash); err != nil || !match {
		return handlerErrors.NewUnauthorizedError("Invalid email or password", InvalidCredentialsCode, requestID, nil)
	}

	// Check if the user has TOTP enabled
	if user.TotpEnabled {
		// Generate a temporary token for TOTP verification
		tokenService := authService.NewTokenService()
		tempToken, err := tokenService.GenerateTempToken(ctx, user.UserID)
		if err != nil {
			return handlerErrors.NewInternalServerError("Failed to generate temporary token", requestID, err)
		}

		// Return the temporary token to be used for TOTP verification
		return c.JSON(http.StatusOK, tempTokenResponse{
			UserID:    user.UserID,
			TempToken: tempToken.TempToken,
			ExpiresIn: tempToken.ExpiresIn,
			Message:   "TOTP verification required",
		})
	}

	// If TOTP is not enabled, generate a token pair directly
	return generateAndReturnTokens(c, ctx, &user, authService.ClientTypeWeb, requestID)
}

// verifyTOTP handles the second step of authentication for users with TOTP enabled
func verifyTOTP(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""

	// Parse and validate request body
	var body totpVerifyRequest
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := middleware.Validate(c, &body); err != nil {
		return err
	}

	// Validate the temporary token
	tokenService := authService.NewTokenService()
	payload, err := tokenService.ValidateToken(ctx, body.TempToken)
	if err != nil {
		return handlerErrors.NewUnauthorizedError("Invalid or expired temporary token", InvalidTOTPTokenCode, requestID, err)
	}

	// Ensure it's a temporary token
	if payload.TokenType != authService.TokenTypeTemp || !payload.RequiresTotp {
		return handlerErrors.NewUnauthorizedError("Invalid token type", InvalidTOTPTokenCode, requestID, nil)
	}

	// Get the user
	q, err := db.Query()
	if err != nil {
		return handlerErrors.NewInternalServerError("Database error", requestID, err)
	}

	user, err := q.GetUserByID(ctx, payload.UserID)
	if err != nil {
		if db.IsNoRows(err) {
			return handlerErrors.NewUnauthorizedError("User not found", UserNotFoundCode, requestID, nil)
		}
		return handlerErrors.NewInternalServerError("Failed to retrieve user", requestID, err)
	}

	// Verify the TOTP code
	if user.TotpSecret == nil {
		return handlerErrors.NewInternalServerError("User TOTP is not configured properly", requestID, nil)
	}

	// Verify TOTP code using a TOTP library (e.g., pquerna/otp)
	// For now, we'll simulate this check
	if body.TOTPCode != "123456" { // Replace with actual TOTP verification logic
		return handlerErrors.NewUnauthorizedError("Invalid TOTP code", InvalidTOTPCodeCode, requestID, nil)
	}

	// TOTP verification successful, generate tokens
	return generateAndReturnTokens(c, ctx, &user, authService.ClientTypeWeb, requestID)
}

// cliLogin handles API-based login for CLI clients
func cliLogin(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""

	// Parse and validate request body
	var body signinRequest
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := middleware.Validate(c, &body); err != nil {
		return err
	}

	// Get the user by email
	q, err := db.Query()
	if err != nil {
		return handlerErrors.NewInternalServerError("Database error", requestID, err)
	}

	user, err := q.GetUserByEmail(ctx, body.Email)
	if err != nil {
		if db.IsNoRows(err) {
			return handlerErrors.NewUnauthorizedError("Invalid email or password", InvalidCredentialsCode, requestID, nil)
		}
		return handlerErrors.NewInternalServerError("Failed to retrieve user", requestID, err)
	}

	// Verify the password
	if match, err := utils.VerifyPassword(body.Password, user.PasswordHash); err != nil || !match {
		return handlerErrors.NewUnauthorizedError("Invalid email or password", InvalidCredentialsCode, requestID, err)
	}

	// Check if the user has TOTP enabled
	if user.TotpEnabled {
		// Return error indicating TOTP verification is required
		// CLI should prompt for TOTP code and call verifyCliTOTP
		return handlerErrors.NewUnauthorizedError("TOTP verification required", RequiresTotpCode, requestID, nil)
	}

	// If TOTP is not enabled, generate a token pair for CLI
	return generateAndReturnTokens(c, ctx, &user, authService.ClientTypeCLI, requestID)
}

// verifyCliTOTP handles TOTP verification for CLI clients
func verifyCliTOTP(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""

	// Parse and validate request body
	var body cliTotpVerifyRequest
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := middleware.Validate(c, &body); err != nil {
		return err
	}

	// Get the user by email
	q, err := db.Query()
	if err != nil {
		return handlerErrors.NewInternalServerError("Database error", requestID, err)
	}

	user, err := q.GetUserByEmail(ctx, body.Email)
	if err != nil {
		if db.IsNoRows(err) {
			return handlerErrors.NewUnauthorizedError("Invalid email", UserNotFoundCode, requestID, err)
		}
		return handlerErrors.NewInternalServerError("Failed to retrieve user", requestID, err)
	}

	// Verify the password again for security
	if match, err := utils.VerifyPassword(body.Password, user.PasswordHash); err != nil || !match {
		return handlerErrors.NewUnauthorizedError("Invalid password", InvalidCredentialsCode, requestID, err)
	}

	// Verify the TOTP code
	if user.TotpSecret == nil {
		return handlerErrors.NewInternalServerError("User TOTP is not configured properly", requestID, nil)
	}

	// Verify TOTP code using a TOTP library (e.g., pquerna/otp)
	// For now, we'll simulate this check
	if body.TOTPCode != "123456" { // Replace with actual TOTP verification logic
		return handlerErrors.NewUnauthorizedError("Invalid TOTP code", InvalidTOTPCodeCode, requestID, nil)
	}

	// TOTP verification successful, generate tokens for CLI
	return generateAndReturnTokens(c, ctx, &user, authService.ClientTypeCLI, requestID)
}

// refreshTokens handles token refresh requests
func refreshTokens(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""
	tokenService := authService.NewTokenService()

	// Try to get refresh token from request body first (for CLI/API clients)
	var body refreshTokenRequest
	refreshToken := ""
	clientType := authService.ClientTypeWeb // Default to web

	if err := c.Bind(&body); err == nil && body.RefreshToken != "" {
		refreshToken = body.RefreshToken
		clientType = authService.ClientTypeCLI // Assume CLI if token in body
	} else {
		// If not in body, check for cookie (for web clients)
		cookie, err := c.Cookie(RefreshTokenCookieName)
		if err != nil || cookie.Value == "" {
			return handlerErrors.NewUnauthorizedError("Refresh token required", RefreshTokenRequiredCode, requestID, nil)
		}
		refreshToken = cookie.Value
	}

	// Refresh the tokens
	tokenPair, err := tokenService.RefreshTokens(ctx, refreshToken)
	if err != nil {
		return handlerErrors.NewUnauthorizedError("Invalid or expired refresh token", InvalidRefreshTokenCode, requestID, err)
	}

	// Get the user ID from the new token
	accessPayload, err := tokenService.ParseToken(tokenPair.AccessToken)
	if err != nil {
		return handlerErrors.NewInternalServerError("Failed to parse new access token", requestID, err)
	}

	// For web clients, set the new refresh token as a cookie
	if clientType == authService.ClientTypeWeb {
		// Add the cookie settings based on context - HTTP or HTTPS
		secure := c.Scheme() == "https"
		sameSite := http.SameSiteLaxMode
		if secure {
			sameSite = http.SameSiteStrictMode
		}

		// Create the cookie
		cookie := new(http.Cookie)
		cookie.Name = RefreshTokenCookieName
		cookie.Value = tokenPair.RefreshToken
		cookie.Path = RefreshTokenCookiePath
		cookie.MaxAge = CookieMaxAge
		cookie.HttpOnly = true
		cookie.Secure = secure
		cookie.SameSite = sameSite

		// Set the cookie
		c.SetCookie(cookie)

		// Clear refresh token from response - web clients don't need it in JSON
		tokenPair.RefreshToken = ""
	}

	// Return the token response
	return c.JSON(http.StatusOK, refreshTokenResponse{
		UserID:       accessPayload.UserID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken, // Will be empty for web clients
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}

// logout handles user sign-out
func logout(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""
	tokenService := authService.NewTokenService()

	// Try to get refresh token from request body first (for CLI/API clients)
	var body logoutRequest
	refreshToken := ""
	clientType := authService.ClientTypeWeb // Default to web

	if err := c.Bind(&body); err == nil && body.RefreshToken != "" {
		refreshToken = body.RefreshToken
		clientType = authService.ClientTypeCLI // Assume CLI if token in body
	} else {
		// If not in body, check for cookie (for web clients)
		cookie, err := c.Cookie(RefreshTokenCookieName)
		if err != nil || cookie.Value == "" {
			// If no refresh token, we can't properly logout, but we'll still clear cookies
			// This is not an error condition for web clients
			if clientType == authService.ClientTypeWeb {
				clearRefreshTokenCookie(c)
				return c.NoContent(http.StatusOK)
			}
			return handlerErrors.NewUnauthorizedError("Refresh token required", RefreshTokenRequiredCode, requestID, nil)
		}
		refreshToken = cookie.Value
	}

	// Revoke the refresh token
	err := tokenService.RevokeToken(ctx, refreshToken)
	if err != nil {
		// If the token is already invalid, we still want to clear cookies
		// This is not an error condition for logout
	}

	// For web clients, clear the refresh token cookie
	if clientType == authService.ClientTypeWeb {
		clearRefreshTokenCookie(c)
	}

	// Return success
	return c.NoContent(http.StatusOK)
}

// clearRefreshTokenCookie is a helper function to clear the refresh token cookie
func clearRefreshTokenCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = RefreshTokenCookieName
	cookie.Value = ""
	cookie.Path = RefreshTokenCookiePath
	cookie.MaxAge = -1 // Delete the cookie
	cookie.HttpOnly = true
	cookie.Secure = c.Scheme() == "https"

	c.SetCookie(cookie)
}

// createAPIToken generates a new API token for a user
func createAPIToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""

	// Get authenticated user from context
	user, err := middleware.GetUserFromContext(c.Request().Context())
	if err != nil {
		return handlerErrors.NewUnauthorizedError("Authentication required", AuthenticationRequiredCode, requestID, err)
	}

	// Generate API token
	tokenService := authService.NewTokenService()
	apiToken, err := tokenService.GenerateAPIToken(ctx, user.UserID)
	if err != nil {
		return handlerErrors.NewInternalServerError("Failed to generate API token", requestID, err)
	}

	// Return the API token
	return c.JSON(http.StatusOK, apiTokenResponse{
		APIToken: apiToken,
	})
}

// revokeAPIToken revokes a specific API token
func revokeAPIToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), defaultTimeout)
	defer cancel()

	requestID := ""

	// Get authenticated user from context
	user, err := middleware.GetUserFromContext(c.Request().Context())
	if err != nil {
		return handlerErrors.NewUnauthorizedError("Authentication required", AuthenticationRequiredCode, requestID, err)
	}

	// Parse and validate request body
	var body revokeAPITokenRequest
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := middleware.Validate(c, &body); err != nil {
		return err
	}

	// Revoke the token
	tokenService := authService.NewTokenService()
	err = tokenService.RevokeToken(ctx, body.APIToken)
	if err != nil {
		// If we can parse the token but it belongs to another user, return unauthorized
		payload, parseErr := tokenService.ParseToken(body.APIToken)
		if parseErr == nil && payload.UserID != user.UserID {
			return handlerErrors.NewUnauthorizedError("Cannot revoke token belonging to another user", UnauthorizedTokenAccessCode, requestID, nil)
		}

		return handlerErrors.NewInternalServerError("Failed to revoke API token", requestID, err)
	}

	// Return success
	return c.NoContent(http.StatusOK)
}
