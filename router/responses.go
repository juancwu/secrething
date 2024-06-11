package router

// apiResponse is a generic json response for all APIs
type apiResponse struct {
	StatusCode int `json:"status_code"`
	Message    any `json:"message"`
}

// loginResponse represents the json body that is sent back when
// a user successfully logs in.
type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// newTokenResponse represents the json body that is sent back when
// a new access token is generated.
type newTokenResponse struct {
	AccessToken string `json:"access_token"`
}
