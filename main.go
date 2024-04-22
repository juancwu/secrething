package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/juancwu/konbini/database"
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
}
