package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Using an numerical type to speed up equality comparisons
type AppEnv uint8

// A list of app environment enums
const (
	APP_ENV_INVALID     AppEnv = 0
	APP_ENV_TESTING     AppEnv = 1
	APP_ENV_DEVELOPMENT AppEnv = 2
	APP_ENV_STAGING     AppEnv = 3
	APP_ENV_PRODUCTION  AppEnv = 4
)

var (
	ErrMissingAppEnv            error = errors.New("APP_ENV environment variable must be set")
	ErrMissingDatabaseUrl       error = errors.New("DATABASE_URL environment variable must be set")
	ErrMissingDatabaseAuthToken error = errors.New("DATABASE_AUTH_TOKEN environment variable must be set")
	ErrMissingBackendUrl        error = errors.New("BACKEND_URL environment variable must be set")
	ErrMissingPort              error = errors.New("PORT environment variable must be set")

	ErrInvalidAppEnv error = errors.New("Invalid value for APP_ENV environment variable")
)

// The server configuration struct. This struct should include all
// the different setups that the server needs. Ideally, just use
// the public methods from this struct instead of accessing the
// fields themselves.
type Config struct {
	env     EnvConfig
	version string
}

type EnvConfig struct {
	databaseUrl       string
	databaseAuthToken string
	backendUrl        string
	port              string
	appEnv            AppEnv
}

// Create a new server configuration. This method reads in required environment
// variables too and it will return an error if any is not set.
func New() (*Config, error) {
	c := &Config{
		version: "development",
	}
	err := c.loadEnvironmentVariables()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Gets the database URL and auth token. The return order is the same (url, token)
func (c *Config) GetDatabaseConfig() (string, string) {
	if c.IsTesting() {
		return c.env.databaseUrl, ""
	}
	return c.env.databaseUrl, c.env.databaseAuthToken
}

// Gets the current backend url value. This value differs based on the environment
// varialbe 'BACKEND_URL'. Different environments should have different values.
func (c *Config) GetBackendUrl() string {
	return c.env.backendUrl
}

// Gets the app environment as a unsigned byte
func (c *Config) GetAppEnvironment() AppEnv {
	return c.env.appEnv
}

// Checks if current app environment is in development mode or not.
func (c *Config) IsDevelopment() bool {
	return c.env.appEnv == APP_ENV_DEVELOPMENT
}

// Checks if current app environment is in testing mode or not.
func (c *Config) IsTesting() bool {
	return c.env.appEnv == APP_ENV_TESTING
}

// Checks if current app environment is in staging mode or not.
func (c *Config) IsStaging() bool {
	return c.env.appEnv == APP_ENV_STAGING
}

// Checks if current app environment is in production mode or not.
func (c *Config) IsProduction() bool {
	return c.env.appEnv == APP_ENV_PRODUCTION
}

// Gets formatted port string. I.E: ":8080"
func (c *Config) GetPort() string {
	return ":" + c.env.port
}

// Gets the unformatted port string. I.E: "8080"
func (c *Config) GetRawPort() string {
	return c.env.port
}

// Gets the current version of the application.
func (c *Config) GetVersion() string {
	return c.version
}

// Load and verify that all required environment variables have been set.
// It will log a warning for missing optional environment variables.
func (c *Config) loadEnvironmentVariables() error {
	// --- start required environment variables ---
	env := os.Getenv("APP_ENV")
	if env == "" {
		return ErrMissingAppEnv
	}
	appEnv, err := c.matchAppEnvStrToEnum(env)
	if err != nil {
		return err
	}
	c.env.appEnv = appEnv

	if c.IsDevelopment() {
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("Failed to load .env file")
		}
	}

	c.env.databaseUrl = os.Getenv("DATABASE_URL")
	if c.env.databaseUrl == "" {
		return ErrMissingDatabaseUrl
	}

	c.env.databaseAuthToken = os.Getenv("DATABASE_AUTH_TOKEN")
	if c.env.databaseAuthToken == "" {
		return ErrMissingDatabaseAuthToken
	}

	c.env.backendUrl = os.Getenv("BACKEND_URL")
	if c.env.backendUrl == "" {
		return ErrMissingBackendUrl
	}

	c.env.port = os.Getenv("PORT")
	if c.env.port == "" {
		return ErrMissingPort
	}

	// --- end required environment variables ---

	return nil
}

// Matches the string representation of app environment. The string representation
// is from the environment varaible 'APP_ENV'. The function will returned an error
// if the string representation is not a valid value.
func (c *Config) matchAppEnvStrToEnum(appEnv string) (AppEnv, error) {
	switch appEnv {
	case "testing":
		return APP_ENV_TESTING, nil
	case "development":
		return APP_ENV_DEVELOPMENT, nil
	case "staging":
		return APP_ENV_STAGING, nil
	case "production":
		return APP_ENV_PRODUCTION, nil
	}
	return 0, ErrInvalidAppEnv
}
