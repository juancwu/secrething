package api

type LoginRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required"`
	TOTPCode *string `json:"totp_code,omitempty" validate:"omitnil,omitempty,required,len=6|len=32"`
}

type CheckAuthTokenRequest struct {
	AuthToken string `json:"auth_token" validate:"required"`
}

type SetupTOTPLockRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}
