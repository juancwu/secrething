package db

import (
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/go-libsql"
)

type DBConnector struct {
	url   string
	token string
}

// NewConnector creates a database connector instance that can be use to create connections
// to the database.
func NewConnector(dbUrl string, dbAuthToken string) *DBConnector {
	var dbString string
	if dbAuthToken != "" {
		dbString = fmt.Sprintf("%s?authToken=%s", dbUrl, dbAuthToken)
	} else {
		dbString = dbUrl
	}
	return &DBConnector{url: dbString, token: dbAuthToken}
}

// Connect creates a new connection to the database.
func (c *DBConnector) Connect() (*sql.DB, error) {
	return sql.Open("libsql", c.url)
}
