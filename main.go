package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/juancwu/konbini/database"
	"github.com/labstack/echo/v4"
)

type AuthReqBody struct {
	Email     string `json:"email" validate:"required"`
	Challenge string `json:"challenge" validate:"required"`
}

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
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading env: %v\n", err)
		}
	}

	fmt.Println("Konbini!")
	db := database.New()
	db.Migrate()

	e := echo.New()
	e.Validator = &ReqValidator{validator: validator.New()}

	e.GET("/health", func(c echo.Context) error {
		report := make(map[string]string)

		sqlDB, err := db.Conn.DB()
		if err != nil {
			report["database"] = "down"
		} else if err := sqlDB.Ping(); err != nil {
			report["database"] = "down"
		} else {
			report["database"] = "up"
		}

		return c.JSON(http.StatusOK, report)
	})

	// auth handling
	e.POST("/auth", func(c echo.Context) error {
		auth := new(AuthReqBody)

		// bind the incoming request data
		if err := c.Bind(auth); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(auth); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, auth)
	})

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
