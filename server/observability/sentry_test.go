package observability

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func TestInitSentry(t *testing.T) {
	// Test initialization with valid config
	t.Run("valid config", func(t *testing.T) {
		config := SentryConfig{
			DSN:              "https://test@test.ingest.sentry.io/test",
			Environment:      "test",
			Release:          "1.0.0",
			Debug:            true,
			SampleRate:       0.5,
			TracesSampleRate: 0.1,
			MaxBreadcrumbs:   100,
			EnableTracing:    true,
			ServerName:       "test-server",
		}

		err := InitSentry(config)
		assert.NoError(t, err)

		// Verify that Sentry client was initialized correctly
		// We can check the current client options
		client := sentry.CurrentHub().Client()
		assert.NotNil(t, client)

		options := client.Options()
		assert.Equal(t, config.DSN, options.Dsn)
		assert.Equal(t, config.Environment, options.Environment)
		assert.Equal(t, config.Release, options.Release)
		assert.Equal(t, config.Debug, options.Debug)
		assert.Equal(t, config.SampleRate, options.SampleRate)
		assert.Equal(t, config.TracesSampleRate, options.TracesSampleRate)
		assert.Equal(t, config.MaxBreadcrumbs, options.MaxBreadcrumbs)
		assert.Equal(t, config.EnableTracing, options.EnableTracing)
		assert.Equal(t, config.ServerName, options.ServerName)
	})

	// Test with empty DSN (common in development/testing)
	t.Run("empty dsn", func(t *testing.T) {
		config := SentryConfig{
			DSN:         "",
			Environment: "development",
		}

		err := InitSentry(config)
		assert.NoError(t, err)

		// Even with empty DSN, should initialize successfully
		client := sentry.CurrentHub().Client()
		assert.NotNil(t, client)
		assert.Empty(t, client.Options().Dsn)
	})
}
