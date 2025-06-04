package config

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

var version string

const (
	// Development environment
	Development string = "development"
	// Production environment
	Production string = "production"
)

// Config holds all application configuration
type Config struct {
	// Database configuration
	DB DatabaseConfig `mapstructure:"db"`

	// Server configuration
	Server ServerConfig `mapstructure:"server"`

	// CORS configuration
	CORS CORSConfig `mapstructure:"cors"`

	// Authentication configuration
	Auth AuthConfig `mapstructure:"auth"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
	AllowHeaders []string `mapstructure:"allow_headers"`
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
	Env     string `mapstructure:"env"`
	Version string
}

type DatabaseConfig struct {
	URL   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

type AuthConfig struct {
	JWT    JWTConfig    `mapstructure:"jwt"`
	Cookie CookieConfig `mapstructure:"cookie"`
}

type JWTConfig struct {
	Secret                    string `mapstructure:"secret"`
	ExpirationMinutes         int    `mapstructure:"expiration_minutes"`
	ExtendedExpirationMinutes int    `mapstructure:"extended_expiration_minutes"`
}

type CookieConfig struct {
	Secure   bool          `mapstructure:"secure"`
	Domain   string        `mapstructure:"domain"`
	Path     string        `mapstructure:"path"`
	SameSite http.SameSite `mapstructure:"same_site"`
	HttpOnly bool          `mapstructure:"http_only"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig(cfgPath string) (*Config, error) {
	cfg := &Config{}
	cfg.Server.Version = version

	viper.SetConfigFile(cfgPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	requiredKeys := []string{
		// Database keys
		"db.url",
		"db.token",

		// Server keys
		"server.address",
		"server.env",

		// CORS keys
		"cors.allow_origins",
		"cors.allow_methods",
		"cors.allow_headers",

		// Auth config
		"auth.jwt.secret",
		"auth.jwt.expiration_minutes",
		"auth.jwt.extended_expiration_minutes",
		"auth.cookie.secure",
		"auth.cookie.domain",
		"auth.cookie.path",
		"auth.cookie.same_site",
		"auth.cookie.http_only",
	}

	notSetList := make([]string, 0)
	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			notSetList = append(notSetList, key)
		}
	}

	if len(notSetList) > 0 {
		return nil, fmt.Errorf("missing %d required configuration key(s): [%s]", len(notSetList), strings.Join(notSetList, ", "))
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	return cfg, nil
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == Development
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Env == Production
}
