package auth

import "github.com/juancwu/go-valkit/v2/validator"

// Request structs for user creation and authentication

type createUserRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,password"`
	Name     *string `json:"name" validate:"omitnil,omitempty,max=50"`
}

func getCreateUserRequestMessages() validator.ValidationMessages {
	msgs := validator.NewValidationMessages()
	msgs.SetMessage("email", "required", "Email is required.")
	msgs.SetMessage("email", "email", "'{2}' is not a valid email.")
	msgs.SetMessage("password", "required", "Password is required.")
	msgs.SetMessage("name", "max", "Name must not be longer than {2} characters.")
	return msgs
}

type signinRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

func getSigninRequestMessages() validator.ValidationMessages {
	msgs := validator.NewValidationMessages()
	msgs.SetMessage("email", "required", "Email is required.")
	msgs.SetMessage("email", "email", "'{2}' is not a valid email.")
	msgs.SetMessage("password", "required", "Password is required.")
	return msgs
}

// Request for TOTP verification after initial login
type totpVerifyRequest struct {
	TempToken string `json:"temp_token" validate:"required"`
	TOTPCode  string `json:"totp_code" validate:"required,len=6,numeric"`
}

func getTotpVerifyRequestMessages() validator.ValidationMessages {
	msgs := validator.NewValidationMessages()
	msgs.SetMessage("temp_token", "required", "Temporary token is required.")
	msgs.SetMessage("totp_code", "required", "TOTP code is required.")
	msgs.SetMessage("totp_code", "len", "TOTP code must be exactly {2} digits.")
	msgs.SetMessage("totp_code", "numeric", "TOTP code must contain only numbers.")
	return msgs
}

// Request for CLI TOTP verification
type cliTotpVerifyRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	TOTPCode string `json:"totp_code" validate:"required,len=6,numeric"`
}

func getCliTotpVerifyRequestMessages() validator.ValidationMessages {
	msgs := validator.NewValidationMessages()
	msgs.SetMessage("email", "required", "Email is required.")
	msgs.SetMessage("email", "email", "'{2}' is not a valid email.")
	msgs.SetMessage("password", "required", "Password is required.")
	msgs.SetMessage("totp_code", "required", "TOTP code is required.")
	msgs.SetMessage("totp_code", "len", "TOTP code must be exactly {2} digits.")
	msgs.SetMessage("totp_code", "numeric", "TOTP code must contain only numbers.")
	return msgs
}

// Token refresh request
type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"omitempty"`
}

// Logout request for CLI/API clients
type logoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"omitempty"`
}

// API token request/response
type revokeAPITokenRequest struct {
	APIToken string `json:"api_token" validate:"required"`
}

func getRevokeAPITokenRequestMessages() validator.ValidationMessages {
	msgs := validator.NewValidationMessages()
	msgs.SetMessage("api_token", "required", "API token is required.")
	return msgs
}
