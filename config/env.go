package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv will try to load .env file in development and
// will check the existance of required environment variables.
func LoadEnv() error {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	// check for required environment variables
	requiredEnvs := []string{
		"DB_URL",
	}
	for _, key := range requiredEnvs {
		val := os.Getenv(key)
		if val == "" {
			return fmt.Errorf("Missing required environmetn variable '%s'", key)
		}
	}

	return nil
}
