package main

import (
	"flag"

	"github.com/juancwu/secrething/internal/api"
	"github.com/juancwu/secrething/internal/config"
)

func main() {
	// Parse flag for custom config path
	cfgPath := flag.String("config", "app.config.yaml", "Configuration file path")

	// Load configuration
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize API and register routes
	apiHandler := api.New(cfg)
	// Start server
	if err := apiHandler.Start(cfg.Server.Address); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
