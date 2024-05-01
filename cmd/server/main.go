package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/router"
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

func loadPublicKey(pemData []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("Failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

func main() {
	/*
	   Environemnt variables are loaded in the env package when it is imported
	*/

	database.Connect()
	database.Migrate()

	e := echo.New()
	// e.Use(middleware.Logger())
	e.Validator = &ReqValidator{validator: validator.New()}

	router.SetupAuthRoutes(e)

	e.GET("/health", func(c echo.Context) error {
		c.Response().WriteHeader(http.StatusOK)
		return nil
	})

	log.Fatal(e.Start(os.Getenv("PORT")))
}
