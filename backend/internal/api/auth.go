package api

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/juancwu/secrething/internal/auth"
	"github.com/juancwu/secrething/internal/config"
	"github.com/juancwu/secrething/internal/db"
	"github.com/labstack/echo/v4"
	"go.jetify.com/typeid"
)

const auth_cookie_name = "auth_token"

// registerAuthRoutes registers all authentication-related routes
func (api *API) registerAuthRoutes() {
	// Group auth routes
	authGroup := api.Echo.Group("/auth")

	// Register routes
	authGroup.POST("/signup", api.handleSignup)
	authGroup.POST("/signin", api.handleSignin)
	authGroup.POST("/signout", api.handleSignout)
}

// SignupRequest represents the payload for user registration
type SignupRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,password"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// SigninRequest represents the payload for user login
type SigninRequest struct {
	Email      string `json:"email" validate:"required,email" errmsg-email:"{value} is not a valid email."`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me" validate:"boolean"`
}

// AuthResponse represents the successful authentication response
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"user"`
	apiResponse
}

// handleSignup handles user registration
func (api *API) handleSignup(c echo.Context) error {
	var req SignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apiResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Validate request
	if err := api.Valkit.Validate(&req); err != nil {
		return err
	}

	// Check if user already exists
	_, err := api.DB.GetUserByEmail(c.Request().Context(), req.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, apiErrorResponse{
			Errors: map[string]string{
				"email": "Email already registered",
			},
		})
	} else if err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check existing user",
		})
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to hash password",
		})
	}

	// Generate user ID
	userID, _ := typeid.New[db.UserID]()

	// Create timestamp
	now := time.Now().UTC().Format(time.RFC3339)

	// Create user
	user, err := api.DB.CreateUser(c.Request().Context(), db.CreateUserParams{
		UserID:       userID,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user, api.Config.Auth.JWT.Secret, api.Config.Auth.JWT.ExpirationMinutes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
		})
	}

	// Set auth cookie
	setAuthCookie(c, token, api.Config, api.Config.Auth.JWT.ExpirationMinutes)

	// Prepare response
	response := AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(api.Config.Auth.JWT.ExpirationMinutes) * time.Minute).Unix(),
	}
	response.Message = "Sign up successful."
	response.User.ID = user.UserID.String()
	response.User.Email = user.Email
	response.User.FirstName = user.FirstName
	response.User.LastName = user.LastName

	return c.JSON(http.StatusCreated, response)
}

// handleSignin handles user login
func (api *API) handleSignin(c echo.Context) error {
	var req SigninRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apiResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Validate request
	if err := api.Valkit.Validate(&req); err != nil {
		return err
	}

	// Get user by email
	user, err := api.DB.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, apiResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid credentials",
			})
		}
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve user",
		})
	}

	// Verify password
	valid, err := auth.CompareHashes(req.Password, user.PasswordHash)
	if err != nil || !valid {
		return c.JSON(http.StatusUnauthorized, apiResponse{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
		})
	}

	// Determine expiration based on remember_me
	expirationMinutes := api.Config.Auth.JWT.ExpirationMinutes
	if req.RememberMe {
		// 30 days for remember me
		expirationMinutes = 60 * 24 * 30
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user, api.Config.Auth.JWT.Secret, expirationMinutes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
		})
	}

	// Create session if remember me is enabled
	if req.RememberMe {
		// Generate session ID
		sessionID, _ := typeid.New[db.SessionID]()

		// Hash the token for storage
		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])

		// Create session
		now := time.Now().UTC()
		expiresAt := now.Add(time.Duration(expirationMinutes) * time.Minute)

		_, err = api.DB.CreateSession(c.Request().Context(), db.CreateSessionParams{
			SessionID:  sessionID,
			UserID:     user.UserID,
			TokenHash:  tokenHash,
			ExpiresAt:  expiresAt.Format(time.RFC3339),
			CreatedAt:  now.Format(time.RFC3339),
			LastUsedAt: now.Format(time.RFC3339),
		})
		if err != nil {
			// Log error but don't fail the signin
			api.Echo.Logger.Errorf("Failed to create session: %v", err)
		}
	}

	// Set auth cookie
	setAuthCookie(c, token, api.Config, expirationMinutes)

	// Prepare response
	response := AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(expirationMinutes) * time.Minute).Unix(),
	}
	response.Message = "Sign in successful."
	response.User.ID = user.UserID.String()
	response.User.Email = user.Email
	response.User.FirstName = user.FirstName
	response.User.LastName = user.LastName

	return c.JSON(http.StatusOK, response)
}

// handleSignout handles user logout
func (api *API) handleSignout(c echo.Context) error {
	// Try to get token to delete session
	var tokenString string

	// Try to get from cookie first
	cookie, err := c.Cookie(auth_cookie_name)
	if err == nil && cookie.Value != "" {
		tokenString = cookie.Value
	}

	// If found, delete the session from database
	if tokenString != "" {
		// Hash the token
		hash := sha256.Sum256([]byte(tokenString))
		tokenHash := hex.EncodeToString(hash[:])

		// Delete session
		err := api.DB.DeleteSessionByTokenHash(c.Request().Context(), tokenHash)
		if err != nil {
			// Log error but don't fail the signout
			api.Echo.Logger.Errorf("Failed to delete session: %v", err)
		}
	}

	// Clear auth cookie by setting it to expire immediately
	unsetAuthCookie(c, api.Config)

	return c.JSON(http.StatusOK, apiResponse{Message: "Signed out successfully"})
}

// Helper function to set auth cookie
func setAuthCookie(c echo.Context, token string, cfg *config.Config, expirationMinutes int) {
	cookie := new(http.Cookie)
	cookie.Name = auth_cookie_name
	cookie.Value = token
	cookie.Path = cfg.Auth.Cookie.Path
	cookie.Domain = cfg.Auth.Cookie.Domain
	cookie.Expires = time.Now().UTC().Add(time.Duration(expirationMinutes) * time.Minute)
	cookie.HttpOnly = cfg.Auth.Cookie.HttpOnly
	cookie.Secure = cfg.Auth.Cookie.Secure
	cookie.SameSite = cfg.Auth.Cookie.SameSite
	c.SetCookie(cookie)
}

// Helper function to unset auth cookie
func unsetAuthCookie(c echo.Context, cfg *config.Config) {
	cookie := new(http.Cookie)
	cookie.Name = auth_cookie_name
	cookie.Value = ""
	cookie.Path = cfg.Auth.Cookie.Path
	cookie.Domain = cfg.Auth.Cookie.Domain
	cookie.Expires = time.Now().Add(-1 * time.Hour) // Set to expire
	cookie.HttpOnly = cfg.Auth.Cookie.HttpOnly
	cookie.Secure = cfg.Auth.Cookie.Secure
	cookie.SameSite = cfg.Auth.Cookie.SameSite
	c.SetCookie(cookie)
}
