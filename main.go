package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/juancwu/konbini/database"
	"github.com/labstack/echo/v4"
)

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

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
