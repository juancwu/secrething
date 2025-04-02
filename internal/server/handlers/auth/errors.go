package auth

// Error codes for authentication
const (
	// Registration errors
	UserAlreadyExistsCode string = "ERR_USR_EMAIL_DUP_4023"
	RegistrationErrorCode string = "ERR_REGISTER_4024"

	// Login errors
	InvalidCredentialsCode     string = "ERR_AUTH_INVALID_CREDS_4010"
	UserNotFoundCode           string = "ERR_AUTH_USER_NOT_FOUND_4011"
	UserAccountLockedCode      string = "ERR_AUTH_ACCOUNT_LOCKED_4012"
	UserEmailNotVerifiedCode   string = "ERR_AUTH_EMAIL_NOT_VERIFIED_4013"
	AuthenticationRequiredCode string = "ERR_AUTH_REQUIRED_4014"

	// TOTP errors
	RequiresTotpCode     string = "ERR_AUTH_TOTP_REQUIRED_4020"
	InvalidTOTPCodeCode  string = "ERR_AUTH_INVALID_TOTP_4021"
	InvalidTOTPTokenCode string = "ERR_AUTH_INVALID_TOTP_TOKEN_4022"

	// Token errors
	InvalidRefreshTokenCode     string = "ERR_AUTH_INVALID_REFRESH_4030"
	RefreshTokenRequiredCode    string = "ERR_AUTH_REFRESH_REQUIRED_4031"
	UnauthorizedTokenAccessCode string = "ERR_AUTH_UNAUTHORIZED_TOKEN_4032"
)
