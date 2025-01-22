package api

type CheckAuthResponse struct {
	AuthToken string `json:"token"`
	TokenType string `json:"type"`
}
