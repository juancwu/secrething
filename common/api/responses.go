package api

type CheckAuthResponse struct {
	AuthToken string `json:"token"`
	TokenType string `json:"type"`
	Email     string `json:"email"`
	// EmailVerified indicates if email has been verified
	EmailVerified bool `json:"email_verified"`
	// TOTP indicates if TOTP has been setup
	TOTP bool `json:"totp"`
}

type RegisterResponse struct {
	AuthToken string `json:"token"`
	TokenType string `json:"type"`
}

type ErrorResponse struct {
	Code      int      `json:"code"`
	Message   string   `json:"message"`
	Errors    []string `json:"errors,omitempty"`
	RequestId string   `json:"request_id"`
}
