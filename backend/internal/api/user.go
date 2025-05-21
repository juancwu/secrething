package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterUserRoutes registers all user-related routes
func (api *API) registerUserRoutes() {
	// Create a group for user routes with auth middleware
	userGroup := api.Echo.Group("/api/user")
	userGroup.Use(api.AuthMiddleware())

	// Protected routes
	userGroup.GET("/profile", api.handleGetProfile)
}

// UserProfile represents the user profile data
type UserProfile struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	apiResponse
}

// handleGetProfile returns the authenticated user's profile
func (api *API) handleGetProfile(c echo.Context) error {
	// Get user from context (set by auth middleware)
	user, ok := c.Get("user").(User)
	if !ok {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user from context",
		})
	}

	// Get user details from database
	dbUser, err := api.DB.GetUserByID(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apiResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve user data",
		})
	}

	// Return user profile
	profile := UserProfile{
		ID:        dbUser.UserID.String(),
		Email:     dbUser.Email,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		apiResponse: apiResponse{
			Code:    http.StatusOK,
			Message: "User profile retrieved successfully",
		},
	}

	return c.JSON(http.StatusOK, profile)
}
