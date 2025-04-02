package responses

// Register defines the response schema for user registration
type Register struct {
	UserID       string  `json:"user_id"`
	Email        string  `json:"email"`
	Name         *string `json:"name,omitempty"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token,omitempty"` // Will be empty for web clients
	ExpiresIn    int     `json:"expires_in"`              // access token lifetime in seconds
}

// Login defines the response schema for user authentication
type Login struct {
	UserID       string  `json:"user_id"`
	Email        string  `json:"email"`
	Name         *string `json:"name,omitempty"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token,omitempty"` // Will be empty for web clients
	ExpiresIn    int     `json:"expires_in"`              // access token lifetime in seconds
}
