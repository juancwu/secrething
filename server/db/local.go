package db

import (
	"database/sql"
	_ "github.com/tursodatabase/go-libsql"
)

// localConnection is a variable to keep the same connection for local since
// the connection persists.
var localConnection *sql.DB

// connectToLocal uses the local compatible libsql library to create a new connection to the database.
func connectToLocal(dbString string) (*sql.DB, error) {
	// re-used the local connection
	if localConnection != nil {
		return localConnection, nil
	}

	localConnection, err := sql.Open("libsql", dbString)
	return localConnection, err
}
