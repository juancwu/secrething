package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

var version string

// Config holds all application configuration
type Config struct {
	// Database configuration
	DB struct {
		URL   string `env:"DB_URL" env-default:"file:./.local/local.db" env-description:"Database connection URL"`
		Token string `env:"DB_TOKEN" env-description:"Authentication token for remote database connections"`
	}

	// Server configuration
	Server struct {
		Address string `env:"SERVER_ADDRESS" env-default:":3000" env-description:"Address and port for the server to listen on"`
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
