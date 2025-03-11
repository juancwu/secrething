package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	DatabaseURL       string
	DatabaseAuthToken string
	BackendURL        string
	Port              string
	Environment       Environment
	Debug             bool
	ServerName        string

	// Email configuration
	ResendAPIKey              string
	VerifyEmailAddress        string
	InvitationEmailAddress    string
	PasswordResetEmailAddress string

	// Security keys
	AuthTokenKey  string
	BentoTokenKey string
	EmailTokenKey string
	AESKey        string

	// Observability
	SentryDSN string

	// Auth settings
	AccessTokenDuration      time.Duration
	EmailVerifyTokenDuration time.Duration
	PasswordResetDuration    time.Duration
	RecoveryCodeCount        int
	RecoveryCodeLength       int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		// Default values
		Port:                     "3000",
		Environment:              EnvDevelopment,
		Debug:                    false,
		AccessTokenDuration:      24 * time.Hour,
		EmailVerifyTokenDuration: 48 * time.Hour,
		PasswordResetDuration:    1 * time.Hour,
		RecoveryCodeCount:        6,
		RecoveryCodeLength:       10,
	}

	// Required environment variables
	requiredVars := []struct {
		name   string
		setter func(string) error
	}{
		{"DATABASE_URL", func(val string) error {
			config.DatabaseURL = val
			return nil
		}},
		{"BACKEND_URL", func(val string) error {
			config.BackendURL = val
			return nil
		}},
		{"SERVER_NAME", func(val string) error {
			config.ServerName = val
			return nil
		}},
		{"AUTH_TOKEN_KEY", func(val string) error {
			if len(val) != 64 {
				return fmt.Errorf("AUTH_TOKEN_KEY must be a 256-bit hexadecimal string (64 characters)")
			}
			config.AuthTokenKey = val
			return nil
		}},
		{"BENTO_TOKEN_KEY", func(val string) error {
			if len(val) != 64 {
				return fmt.Errorf("BENTO_TOKEN_KEY must be a 256-bit hexadecimal string (64 characters)")
			}
			config.BentoTokenKey = val
			return nil
		}},
		{"EMAIL_TOKEN_KEY", func(val string) error {
			if len(val) != 64 {
				return fmt.Errorf("EMAIL_TOKEN_KEY must be a 256-bit hexadecimal string (64 characters)")
			}
			config.EmailTokenKey = val
			return nil
		}},
		{"AES_KEY", func(val string) error {
			if len(val) != 64 {
				return fmt.Errorf("AES_KEY must be a 256-bit hexadecimal string (64 characters)")
			}
			config.AESKey = val
			return nil
		}},
	}

	// Email-related required variables
	emailRequiredVars := []struct {
		name   string
		setter func(string) error
	}{
		{"RESEND_API_KEY", func(val string) error {
			config.ResendAPIKey = val
			return nil
		}},
		{"VERIFY_EMAIL_ADDRESS", func(val string) error {
			config.VerifyEmailAddress = val
			return nil
		}},
		{"INVITATION_EMAIL_ADDRESS", func(val string) error {
			config.InvitationEmailAddress = val
			return nil
		}},
		{"PASSWORD_RESET_EMAIL_ADDRESS", func(val string) error {
			config.PasswordResetEmailAddress = val
			return nil
		}},
	}

	// Validate and set required variables
	for _, v := range requiredVars {
		val := os.Getenv(v.name)
		if val == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", v.name)
		}
		if err := v.setter(val); err != nil {
			return nil, err
		}
	}

	// Email config is optional in development mode
	envStr := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if envStr == "" {
		envStr = string(EnvDevelopment)
	}

	config.Environment = Environment(envStr)

	for _, v := range emailRequiredVars {
		val := os.Getenv(v.name)
		if val == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", v.name)
		}
		if err := v.setter(val); err != nil {
			return nil, err
		}
	}

	// Optional variables with defaults
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}

	if debug := os.Getenv("DEBUG"); debug != "" {
		parsedDebug, err := strconv.ParseBool(debug)
		if err == nil {
			config.Debug = parsedDebug
		}
	}

	if dbAuthToken := os.Getenv("DATABASE_AUTH_TOKEN"); dbAuthToken != "" {
		config.DatabaseAuthToken = dbAuthToken
	}

	if sentryDSN := os.Getenv("SENTRY_DSN"); sentryDSN != "" {
		config.SentryDSN = sentryDSN
	}

	return config, nil
}

// IsDevelopment returns true if the environment is set to development
func (c *Config) IsDevelopment() bool {
	return c.Environment.IsDevelopment()
}

// IsProduction returns true if the environment is set to production
func (c *Config) IsProduction() bool {
	return c.Environment.IsProduction()
}

// IsStaging returns true if the environment is set to staging
func (c *Config) IsStaging() bool {
	return c.Environment.IsStaging()
}

// IsTest returns true if the environment is set to test
func (c *Config) IsTest() bool {
	return c.Environment.IsTest()
}

// GetAddress returns the address the server should listen on
func (c *Config) GetAddress() string {
	return fmt.Sprintf(":%s", c.Port)
}
