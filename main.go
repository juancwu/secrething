package main

import (
	"os"

	// package modules
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	// custom modules
	"github.com/juancwu/konbini/config"
	"github.com/juancwu/konbini/middleware"
	"github.com/juancwu/konbini/router"
	"github.com/juancwu/konbini/store"
)

type customValidator struct {
	validator *validator.Validate
}

func (v *customValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

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

	e := echo.New()
	e.Use(middleware.UseRequestId())
	e.Validator = &customValidator{validator: validator.New()}
	// remove the banner and port logging in production
	e.HideBanner = os.Getenv("APP_ENV") != config.DEV_ENV
	e.HidePort = os.Getenv("APP_ENV") != config.DEV_ENV
	api := e.Group("/api/v1")
	router.SetupHealthcheckRoutes(api)
	router.SetupAccountRoutes(api)

	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		logger.Fatal("Failed to start http server", zap.Error(err))
	}
}
