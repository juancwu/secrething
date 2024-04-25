package database

import (
	"github.com/charmbracelet/log"
	"github.com/juancwu/konbini/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	gormdb, err := gorm.Open(postgres.Open(env.Values.DB_URL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	db = gormdb
}

func GetDB() *gorm.DB {
	return db
}
