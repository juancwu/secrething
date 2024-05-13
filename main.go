package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/juancwu/konbini/server/database"
	_ "github.com/juancwu/konbini/server/env"
	"github.com/juancwu/konbini/server/router"
	"github.com/juancwu/konbini/server/service"
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

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Konbini is healthy")
	})

	e.GET("/verify/token", func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) < 2 {
			return c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}
		token := parts[1]
		utils.Logger().Infof("Token: %s\n", token)
		parsedToken, err := service.VerifyToken(token)
		if err != nil {
			utils.Logger().Errorf("Failed to verify token: %v\n", err)
			return c.String(http.StatusUnauthorized, "failed to verify token")
		}
		if parsedToken.Valid {
			utils.Logger().Infof("Parsed token: %v\n", parsedToken.Claims)
		} else {
			utils.Logger().Info("Invalid token")
			return c.String(http.StatusUnauthorized, "invalid token")
		}
		return c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})

	log.Fatal(e.Start(os.Getenv("PORT")))
}
