package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	// DEV_ENV represents the development environment of the app
	DEV_ENV string = "development"
	// PROD_ENV represents the production environment of the app
	PROD_ENV string = "production"
)

// LoadEnv will try to load .env file in development and
// will check the existance of required environment variables.
func LoadEnv() error {
	if os.Getenv("APP_ENV") == DEV_ENV {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	// check for required environment variables
	requiredEnvs := []string{
		"DB_URL",
		"PORT",
		"RESEND_API_KEY",
		"SERVER_URL",
		"PGP_SYM_KEY",
		"PASS_ENCRYPT_ALGO",
		"NOREPLY_EMAIL",
		"DONOTREPLY_EMAIL",
		"JWT_ACCESS_TOKEN_TYPE",
		"JWT_REFRESH_TOKEN_TYPE",
		"JWT_ACCESS_TOKEN_SECRET",
		"JWT_REFRESH_TOKEN_SECRET",
		"JWT_ISSUER",
	}
	for _, key := range requiredEnvs {
		val := os.Getenv(key)
		if val == "" {
			return fmt.Errorf("Missing required environmetn variable '%s'", key)
		}
	}

	return nil
}
