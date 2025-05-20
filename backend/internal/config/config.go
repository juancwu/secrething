package config

import (
	"fmt"
	"net/http"

	"github.com/ilyakaznacheev/cleanenv"
)

var version string

// Environment represents the runtime environment
type Environment string

const (
	// Development environment
	Development Environment = "development"
	// Production environment
	Production Environment = "production"
)

// Config holds all application configuration
type Config struct {
	// Environment configuration
	Env Environment `env:"APP_ENV" env-default:"development" env-description:"Application environment (development or production)"`

	// Database configuration
	DB struct {
		URL   string `env:"DB_URL" env-default:"file:./.local/local.db" env-description:"Database connection URL"`
		Token string `env:"DB_TOKEN" env-description:"Authentication token for remote database connections"`
	}

	// Server configuration
	Server struct {
		Address string `env:"SERVER_ADDRESS" env-default:":3000" env-description:"Address and port for the server to listen on"`
	}

	// CORS configuration
	CORS struct {
		AllowOrigins []string `env:"CORS_ALLOW_ORIGINS" env-default:"http://localhost:5173,https://secrething.app" env-description:"Comma-separated list of allowed origins"`
		AllowMethods []string `env:"CORS_ALLOW_METHODS" env-default:"GET,POST,PUT,DELETE,OPTIONS" env-description:"Comma-separated list of allowed HTTP methods"`
		AllowHeaders []string `env:"CORS_ALLOW_HEADERS" env-default:"Accept,Authorization,Content-Type,X-CSRF-Token" env-description:"Comma-separated list of allowed HTTP headers"`
	}

	// Authentication configuration
	Auth struct {
		JWTSecret            string        `env:"JWT_SECRET" env-description:"Secret key for JWT token generation"`
		JWTExpirationMinutes int           `env:"JWT_EXPIRATION_MINUTES" env-default:"60" env-description:"JWT token expiration time in minutes"`
		CookieDomain         string        `env:"COOKIE_DOMAIN" env-description:"Domain for auth cookies"`
		CookieSecure         bool          `env:"COOKIE_SECURE" env-default:"false" env-description:"Whether cookies should only be sent over HTTPS"`
		CookiePath           string        `env:"COOKIE_PATH" env-default:"" env-description:"Path cookie is included in requests"`
		CookieSameSite       http.SameSite `env:"COOKIE_SAME_SITE" env-default:"3" env-description:"Same site mode for cookie. Default: strict(3)"`
		CookieHttpOnly       bool          `env:"COOKIE_HTTP_ONLY" env-default:"true" env-description:"Set cookie to be http-only"`
	}

	// Application configuration
	App struct {
		Version string `env:"-"`
	}
}

// NewConfig creates a new Config instance with default values
func NewConfig() *Config {
	cfg := &Config{}
	cfg.App.Version = version
	return cfg
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := NewConfig()

	// Read from environment variables (.env file loaded separately)
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}

	return cfg, nil
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == Development
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.Env == Production
}
