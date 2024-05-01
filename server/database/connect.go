package database

import (
	"database/sql"

	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"

	"github.com/juancwu/konbini/server/env"
)

var db *sql.DB

func Connect() {
	var err error
	db, err = sql.Open("postgres", env.Values().DB_URL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
}

func DB() *sql.DB {
	return db
}
