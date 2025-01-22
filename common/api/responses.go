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
