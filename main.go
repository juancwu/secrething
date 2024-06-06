package main

import (
	// package modules
	"os"

	"go.uber.org/zap"

	// custom modules
	"github.com/juancwu/konbini/config"
	"github.com/juancwu/konbini/store"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		logger.Fatal("Failed to load env", zap.Error(err))
	}

	err = store.Connect(os.Getenv("DB_URL"))
	if err != nil {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		logger.Fatal("Failed to establish connection with database", zap.Error(err))
	}
}
