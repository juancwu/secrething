package test

import (
	"konbini/server/config"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
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
		require.True(t, c.IsTesting())
	})

	t.Run("correct development app environment", func(t *testing.T) {
		original := os.Getenv("APP_ENV")
		os.Setenv("APP_ENV", "development")
		defer func(value string) { os.Setenv("APP_ENV", value) }(original)

		c, err := config.New()
		require.NoError(t, err)

		require.True(t, c.IsDevelopment())

		_, token := c.GetDatabaseConfig()
		require.Equal(t, os.Getenv("DATABASE_AUTH_TOKEN"), token)
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
