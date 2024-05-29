package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/juancwu/konbini/server/database"
	_ "github.com/juancwu/konbini/server/env"
	"github.com/juancwu/konbini/server/router"
	"github.com/juancwu/konbini/server/utils"
)

type ReqValidator struct {
	validator *validator.Validate
}

func (rq *ReqValidator) Validate(i interface{}) error {
	if err := rq.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func main() {
	/*
	   Environemnt variables are loaded in the env package when it is imported
	*/

	database.Connect()
	database.Migrate()

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:        true,
		LogRoutePath:     true,
		LogMethod:        true,
		LogError:         true,
		LogRemoteIP:      true,
		LogUserAgent:     true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				utils.Logger().Info(fmt.Sprintf("[%s Request]", v.Method), "route", v.RoutePath, "status", v.Status, "remote-ip", v.RemoteIP, "agent", v.UserAgent, "content-length", v.ContentLength, "res-size", v.ResponseSize)
			} else {
				utils.Logger().Error(fmt.Sprintf("[%s Request]", v.Method), "route", v.RoutePath, "status", v.Status, "remote-ip", v.RemoteIP, "agent", v.UserAgent, "content-length", v.ContentLength, "res-size", v.ResponseSize, "error", v.Error)
			}
			return nil
		},
	}))
	e.Validator = &ReqValidator{validator: validator.New()}

	router.SetupAuthRoutes(e)
	router.SetupBentoRoutes(e)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Konbini is healthy (%s)", os.Getenv("APP_VERSION")))
	})

	log.Fatal(e.Start(os.Getenv("PORT")))
}
