package config

// DatabaseConfig represents the configuration settings for database connections.
// It contains all parameters needed to establish and maintain connections to the database.
type DatabaseConfig struct {
	// URL specifies the database connection string
	// For SQLite: "file:/path/to/db.sqlite"
	// For remote DBs like Turso: protocol://hostname:port/db
	URL string `env:"DATABASE_URL" env-default:"file:./.local/local.db" env-required:"" env-description:"The URL to connect to the database."`

	// AuthToken holds the authentication token for database services
	// that require token-based auth (like Turso)
	AuthToken string `env:"DATABASE_AUTH_TOKEN" env-default:"" env-required:"" env-description:"The auth token to connect to the database."`
}

// databaseCfg is the private singleton instance of database configuration.
// It's populated when Load() is called.
var databaseCfg DatabaseConfig

// Database returns a copy of the initialized database configuration.
// The configuration must be loaded with Load() before this function is called.
func Database() DatabaseConfig {
	return databaseCfg
}
