package test

import (
	"konbini/server/config"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("uninitialized global configuration", func(t *testing.T) {
		c, err := config.Global()
		require.Nil(t, c)
		require.NotNil(t, err)
		require.ErrorIs(t, err, config.ErrUninitializedGlobalConfig)
	})

	t.Run("new server configuration", func(t *testing.T) {
		c, err := config.New()
		require.NoError(t, err)

		url, token := c.GetDatabaseConfig()
		require.Equal(t, os.Getenv("DATABASE_URL"), url)
		// expecting empty token string for testing environment
		require.Equal(t, "", token)
		require.Equal(t, os.Getenv("BACKEND_URL"), c.GetBackendUrl())
		require.Equal(t, os.Getenv("PORT"), c.GetRawPort())
		require.Equal(t, ":"+os.Getenv("PORT"), c.GetPort())
		require.Equal(t, os.Getenv("RESEND_API_KEY"), c.GetResendApiKey())
		require.Equal(t, os.Getenv("NOREPLY_EMAIL"), c.GetNoReplyEmail())
		require.Equal(t, os.Getenv("BENTO_TOKEN_ISSUER"), c.GetBentoTokenIssuer())
		require.Equal(t, os.Getenv("EMAIL_TOKEN_ISSUER"), c.GetEmailTokenIssuer())
		require.True(t, c.IsTesting())
	})

	t.Run("global config", func(t *testing.T) {
		c, err := config.Global()
		require.NoError(t, err)
		require.NotNil(t, c)
	})

	t.Run("correct development app environment", func(t *testing.T) {
		original := os.Getenv("APP_ENV")
		os.Setenv("APP_ENV", "development")
		defer func(value string) { os.Setenv("APP_ENV", value) }(original)

		// need to create a temporary .env file since development expects
		// to have a .env file in the current working directory
		envPath := ".env"
		backupPath := ".env.bk"

		_, err := os.Stat(envPath)
		hasEnvFile := !os.IsNotExist(err)
		if hasEnvFile {
			if err := os.Rename(envPath, backupPath); err != nil {
				t.Fatalf("Failed to rename original .env file during config development mode test: %v", err)
			}
		}

		if err := os.WriteFile(envPath, []byte(""), 0644); err != nil {
			if hasEnvFile {
				err := os.Rename(backupPath, envPath)
				if err != nil {
					t.Fatalf("Failed to rename backup .env file back to original during config development mode test: %v", err)
				}
			}
			t.Fatalf("Failed to create temporary .env file during config development mode test: %v", err)
		}

		// --- actual test

		c, err := config.New()
		require.NoError(t, err)

		require.True(t, c.IsDevelopment())

		_, token := c.GetDatabaseConfig()
		require.Equal(t, os.Getenv("DATABASE_AUTH_TOKEN"), token)

		// --- end of actual test

		if err := os.Remove(envPath); err != nil {
			t.Fatalf("Failed to remove temporary .env file during development mode test: %v", err)
		}
		if hasEnvFile {
			if err := os.Rename(backupPath, envPath); err != nil {
				t.Fatalf("Failed to rename temporary .env file back during development mode test: %v", err)
			}
		}
	})

	t.Run("correct staging app environment", func(t *testing.T) {
		original := os.Getenv("APP_ENV")
		os.Setenv("APP_ENV", "staging")
		defer func(value string) { os.Setenv("APP_ENV", value) }(original)

		c, err := config.New()
		require.NoError(t, err)

		require.True(t, c.IsStaging())

		_, token := c.GetDatabaseConfig()
		require.Equal(t, os.Getenv("DATABASE_AUTH_TOKEN"), token)
	})

	t.Run("correct production app environment", func(t *testing.T) {
		original := os.Getenv("APP_ENV")
		os.Setenv("APP_ENV", "production")
		defer func(value string) { os.Setenv("APP_ENV", value) }(original)

		c, err := config.New()
		require.NoError(t, err)

		require.True(t, c.IsProduction())

		_, token := c.GetDatabaseConfig()
		require.Equal(t, os.Getenv("DATABASE_AUTH_TOKEN"), token)
	})
}
