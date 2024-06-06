package store

import (
	// builtin modules
	"database/sql"

	// package modules
	_ "github.com/lib/pq"
)

// db is the connection we have established with the database when Connect is successful.
var db *sql.DB

// Connect establishes a connection with the given postgresql database with the given url.
func Connect(dbUrl string) error {
	var err error
	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	return nil
}
