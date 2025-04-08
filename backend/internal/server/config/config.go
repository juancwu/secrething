// Package config provides a centralized configuration system for the Konbini application.
// It handles loading settings from environment variables and .env files,
// providing type-safe access to configuration values through specialized
// configuration objects for different parts of the system.
package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var version string

// Load loads configuration values from environment variables and optional .env files.
// It should be called once during application initialization.
//
// If filenames are provided, it will attempt to load environment variables from those files.
// After loading variables from files (if any), it reads environment variables into
// the specialized configuration structs for different parts of the system.
//
// Returns an error if loading from files fails or if reading environment variables
// into any configuration struct fails.
func Load(filenames ...string) error {
	// First try to load variables from .env files if provided
	if err := godotenv.Load(filenames...); err != nil {
		return err
	}

	// Read environment variables into each configuration struct
	if err := cleanenv.ReadEnv(&serverCfg); err != nil {
		return err
	}
	if err := cleanenv.ReadEnv(&databaseCfg); err != nil {
		return err
	}
	if err := cleanenv.ReadEnv(&emailCfg); err != nil {
		return err
	}
	if err := cleanenv.ReadEnv(&tokenCfg); err != nil {
		return err
	}
	if err := cleanenv.ReadEnv(&cryptoCfg); err != nil {
		return err
	}

	return nil
}

// Version returns the current version of the server.
// This is set during the build process.
func Version() string {
	return version
}
