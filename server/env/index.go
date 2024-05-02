package env

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

type Env struct {
	APP_ENV           string
	PORT              string
	DB_URL            string
	DB_NAME           string
	RESEND_API_KEY    string
	SERVER_URL        string
	PGP_SYM_KEY       string
	PASS_ENCRYPT_ALGO string
	NOREPLY_EMAIL     string
}

var values Env

func init() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading env: %v\n", err)
		}
	}

	values = Env{}

	// required env
	values.PORT = getEnv("PORT", true)
	values.DB_URL = getEnv("DB_URL", true)
	values.DB_NAME = getEnv("DB_NAME", true)
	values.RESEND_API_KEY = getEnv("RESEND_API_KEY", true)
	values.SERVER_URL = getEnv("SERVER_URL", true)
	values.PGP_SYM_KEY = getEnv("PGP_SYM_KEY", true)
	values.PASS_ENCRYPT_ALGO = getEnv("PASS_ENCRYPT_ALGO", true)
	values.NOREPLY_EMAIL = getEnv("NOREPLY_EMAIL", true)

	// optional env
	values.APP_ENV = getEnv("APP_ENV", false)
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

func Values() *Env {
	return &values
}
