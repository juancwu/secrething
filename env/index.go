package env

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

type Env struct {
	APP_ENV string
	PORT    string
	DB_URL  string
}

var Values Env

func init() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading env: %v\n", err)
		}
	}

	Values = Env{}

	// required env
	Values.PORT = getEnv("PORT", true)
	Values.DB_URL = getEnv("DB_URL", true)

	// optional env
	Values.APP_ENV = getEnv("APP_ENV", false)
}

// checks if env exists or not
func getEnv(key string, required bool) string {
	v := os.Getenv(key)
	if v == "" {
		if required {
			log.Fatalf("Missing required env: %s\n", key)
		} else {
			log.Warnf("Missing optional env: %s\n", key)
		}
	}
	return v
}
