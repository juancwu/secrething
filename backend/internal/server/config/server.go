package config

// Server environment constants define the supported runtime environments for the application.
// These constants are used to determine environment-specific behavior and logging.
const (
	// ServerEnvDevelopment represents the local development environment
	ServerEnvDevelopment string = "development"

	// ServerEnvSandbox represents a testing sandbox environment
	ServerEnvSandbox string = "sandbox"

	// ServerEnvStaging represents the pre-production/QA environment
	ServerEnvStaging string = "staging"

	// ServerEnvProduction represents the live production environment
	ServerEnvProduction string = "production"
)

// ServerConfig represents the core server configuration.
// It contains settings that are needed for the server to boot and operate properly.
type ServerConfig struct {
	// Name identifies this server instance, useful in multi-server deployments
	Name string `env:"SERVER_NAME" env-default:"konbini-local" env-description:"The server name. This is to easily identify the different duplicates of the server."`

	// Address specifies the network address and port the server should listen on
	Address string `env:"SERVER_ADDRESS" env-default:":3000" env-required:"" env-description:"The address the server is listening on."`

	// URL is the publicly accessible base URL for this server instance
	URL string `env:"SERVER_URL" env-default:"http://127.0.0.1:3000" env-required:"" env-description:"The URL use to reach the server."`

	// Env specifies the runtime environment (development, sandbox, staging, or production)
	Env string `env:"SERVER_ENV" env-default:"development" env-required:"" env-description:"The server runtime environment."`
}

// serverCfg is the internal singleton instance of the server configuration
var serverCfg ServerConfig

// Server returns a copy of the initialized server configuration.
func Server() ServerConfig {
	return serverCfg
}

// IsDevelopment checks if the current runtime environment is development.
// Use this for enabling development-only features like detailed logging.
func IsDevelopment() bool {
	return serverCfg.Env == ServerEnvDevelopment
}

// IsSandbox checks if the current runtime environment is sandbox.
// The sandbox environment is typically used for testing with external systems.
func IsSandbox() bool {
	return serverCfg.Env == ServerEnvSandbox
}

// IsStaging checks if the current runtime environment is staging.
// The staging environment mimics production but is used for final testing.
func IsStaging() bool {
	return serverCfg.Env == ServerEnvStaging
}

// IsProduction checks if the current runtime environment is production.
// Use this to enable production-specific behavior like stricter security.
func IsProduction() bool {
	return serverCfg.Env == ServerEnvProduction
}
