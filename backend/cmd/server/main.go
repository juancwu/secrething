package main

import (
	"github.com/joho/godotenv"
	"github.com/juancwu/secrething/internal/api"
	"github.com/juancwu/secrething/internal/config"
)

func main() {
	// Load .env file first
	godotenv.Load()

	// Load configuration
	cfg, err := config.LoadConfig()
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
