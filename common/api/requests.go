package api

type CheckAuthTokenRequest struct {
	AuthToken string `json:"auth_token" validate:"required"`
}
