package observability

import "github.com/getsentry/sentry-go"

// SentryConfig holds sentry configuration options
type SentryConfig struct {
	DSN              string
	Environment      string
	Release          string
	Debug            bool
	SampleRate       float64
	TracesSampleRate float64
	MaxBreadcrumbs   int
	EnableTracing    bool
	ServerName       string
}

// InitSentry initializes the Sentry SDK
func InitSentry(config SentryConfig) error {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.DSN,
		Environment:      config.Environment,
		Release:          config.Release,
		Debug:            config.Debug,
		SampleRate:       config.SampleRate,
		TracesSampleRate: config.TracesSampleRate,
		MaxBreadcrumbs:   config.MaxBreadcrumbs,
		EnableTracing:    config.EnableTracing,
		ServerName:       config.ServerName,
	}); err != nil {
		return err
	}

	return nil
}
