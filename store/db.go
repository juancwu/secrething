package store

import (
	"database/sql"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Establishes a new DB connection.
func NewConn() (*sql.DB, error) {
	url := os.Getenv("TURSO_DATABASE_URL")
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
