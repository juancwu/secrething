package config

// EmailConfig holds configuration settings for email services.
// It includes API keys and sender addresses for various types of system emails.
type EmailConfig struct {
	// ResendApiKey is the API key for the Resend.com email service.
	// Used to authenticate requests to the email delivery service.
	ResendApiKey string `env:"EMAIL_RESEND_API_KEY" env-required:"" env-description:"Resend API key use to send emails."`

	// VerifyAddress is the sender email address for account verification emails.
	// This should be a verified domain in your email service.
	VerifyAddress string `env:"EMAIL_VERIFY_ADDRESS" env-required:"" env-description:"The email address use for sending verification emails."`

	// InvitationAddress is the sender email address for team/group invitation emails.
	// Users will see this as the from address when receiving invitations.
	InvitationAddress string `env:"EMAIL_INVITATION_ADDRESS" env-required:"" env-description:"The email address use for sending team/group invitation emails."`

	// PasswordResetAddress is the sender email address for password reset emails.
	// Using a distinct address helps users identify the email purpose.
	PasswordResetAddress string `env:"EMAIL_PASSWORD_RESET_ADDRESS" env-required:"" env-description:"The email address use for sending password reset emails."`
}

// emailCfg is the private singleton instance of email configuration.
// It's populated when Load() is called.
var emailCfg EmailConfig

// Email returns a copy of the initialized email configuration.
// The configuration must be loaded with Load() before this function is called.
func Email() EmailConfig {
	return emailCfg
}
