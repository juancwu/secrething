package database

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/juancwu/konbini/server/env"
	"github.com/juancwu/konbini/server/utils"
)

var db *sql.DB

func Connect() {
	var err error
	db, err = sql.Open("postgres", env.Values().DB_URL)
	if err != nil {
		utils.Logger().Fatalf("Error connecting to database: %v\n", err)
	}
}

func DB() *sql.DB {
	return db
}
