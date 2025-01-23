package api

// common place for routes

const (
	UriLogin                   = "/auth/login"
	UriRegister                = "/auth/register"
	UriCheckToken              = "/auth/token/check"
	UriTOTPSetup               = "/auth/totp/setup"
	UriTOTPLock                = "/auth/totp/lock"
	UriTOTPDelete              = "/auth/totp"
	UriVerifyEmail             = "/auth/email/verify"
	UriResendVerificationEmail = "/auth/email/resend-verification"
)
