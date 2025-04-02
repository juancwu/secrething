package auth

// createUserResponse holds info for the create user endpoint response.
type createUserResponse struct {
	UserID string `json:"user_id"`
}

// loginResponse is returned after successful authentication
type loginResponse struct {
	UserID       string  `json:"user_id"`
	Email        string  `json:"email"`
	Name         *string `json:"name,omitempty"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token,omitempty"` // Will be empty for web clients
	ExpiresIn    int     `json:"expires_in"`              // access token lifetime in seconds
}

// tempTokenResponse is returned when TOTP verification is required
type tempTokenResponse struct {
	UserID    string `json:"user_id"`
	TempToken string `json:"temp_token"`
	ExpiresIn int    `json:"expires_in"` // temporary token lifetime in seconds
	Message   string `json:"message"`
}

// refreshTokenResponse is returned after a token refresh
type refreshTokenResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"` // Will be empty for web clients
	ExpiresIn    int    `json:"expires_in"`              // access token lifetime in seconds
}

// apiTokenResponse is returned when creating a new API token
type apiTokenResponse struct {
	APIToken string `json:"api_token"`
}
