package db

import (
	"database/sql"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// connectToRemote uses the remote compatible libsql package
func connectToRemote(dbString string) (*sql.DB, error) {
	return sql.Open("libsql", dbString)
}
