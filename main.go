package main

import (
	// package modules
	"go.uber.org/zap"

	// custom modules
	"github.com/juancwu/konbini/config"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		logger.Fatal("Failed to load env", zap.Error(err))
	}
}
