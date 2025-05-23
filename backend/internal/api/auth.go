package api

import (
	"database/sql"
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
	Email    string `json:"email" validate:"required,email" errmsg-email:"{value} is not a valid email."`
	Password string `json:"password" validate:"required"`
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
	token, err := auth.GenerateToken(user, api.Config.Auth.JWTSecret, api.Config.Auth.JWTExpirationMinutes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
		})
	}

	// Set auth cookie
	setAuthCookie(c, token, api.Config)

	// Prepare response
	response := AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(api.Config.Auth.JWTExpirationMinutes) * time.Minute).Unix(),
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

	// Generate JWT token
	token, err := auth.GenerateToken(user, api.Config.Auth.JWTSecret, api.Config.Auth.JWTExpirationMinutes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
		})
	}

	// Set auth cookie
	setAuthCookie(c, token, api.Config)

	// Prepare response
	response := AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(api.Config.Auth.JWTExpirationMinutes) * time.Minute).Unix(),
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
	// Clear auth cookie by setting it to expire immediately
	cookie := new(http.Cookie)
	cookie.Name = auth_cookie_name
	cookie.Value = ""
	cookie.Path = api.Config.Auth.CookiePath
	cookie.Domain = api.Config.Auth.CookieDomain
	cookie.Expires = time.Now().Add(-1 * time.Hour) // Set to expire
	cookie.HttpOnly = api.Config.Auth.CookieHttpOnly
	cookie.Secure = api.Config.Auth.CookieSecure
	cookie.SameSite = api.Config.Auth.CookieSameSite
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, apiResponse{Message: "Signed out successfully"})
}

// Helper function to set auth cookie
func setAuthCookie(c echo.Context, token string, cfg *config.Config) {
	cookie := new(http.Cookie)
	cookie.Name = auth_cookie_name
	cookie.Value = token
	cookie.Path = cfg.Auth.CookiePath
	cookie.Domain = cfg.Auth.CookieDomain
	cookie.Expires = time.Now().UTC().Add(time.Duration(cfg.Auth.JWTExpirationMinutes) * time.Minute)
	cookie.HttpOnly = cfg.Auth.CookieHttpOnly
	cookie.Secure = cfg.Auth.CookieSecure
	cookie.SameSite = cfg.Auth.CookieSameSite
	c.SetCookie(cookie)
}
