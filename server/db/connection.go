package db

import (
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/go-libsql"
)

// NewConnection opens a new database connection with the given database url and auth token.
// If local connection is desired, pass an empty string for the auth token parameter.
func NewConnection(dbUrl string, dbAuthToken string) (*sql.DB, error) {
	var dbString string
	if dbAuthToken != "" {
		dbString = fmt.Sprintf("%s?authToken=%s", dbUrl, dbAuthToken)
	} else {
		dbString = dbUrl
	}

	db, err := sql.Open("libsql", dbString)
	return db, err
}
