package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Establishes a new DB connection.
func NewConn() (*sql.DB, error) {
	url := fmt.Sprintf(
		os.Getenv("TURSO_DATABASE_URL")+"?authToken=%s",
		os.Getenv("TURSO_AUTH_TOKEN"),
	)
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
