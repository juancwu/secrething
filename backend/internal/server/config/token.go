package config

// TokenConfig holds encryption keys for various types of security tokens.
// Each token type has its own dedicated encryption key for better security isolation.
type TokenConfig struct {
	// AuthKey is the encryption key used for authentication tokens.
	// These tokens are used for user sessions and API authentication.
	AuthKey string `env:"TOKEN_AUTH_KEY" env-required:"" env-description:"The AES key use to encrypt auth tokens."`

	// BentoKey is the encryption key used for bento (secret container) tokens.
	// These tokens secure access to user secrets and shared containers.
	BentoKey string `env:"TOKEN_BENTO_KEY" env-required:"" env-description:"The AES key use to encrypt bento tokens."`

	// EmailKey is the encryption key used for email verification tokens.
	// These tokens are embedded in verification links sent to users.
	EmailKey string `env:"TOKEN_EMAIL_KEY" env-required:"" env-description:"The AES key use to encrypt email tokens."`
}

// tokenCfg is the private singleton instance of token configuration.
// It's populated when Load() is called.
var tokenCfg TokenConfig

// Token returns a copy of the initialized token configuration.
// The configuration must be loaded with Load() before this function is called.
func Token() TokenConfig {
	return tokenCfg
}
