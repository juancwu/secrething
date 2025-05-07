package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/tursodatabase/go-libsql"
)

var connection *sql.DB
var mu sync.RWMutex

// Connect establishes a connection using the active driver.
// It's safe to call Connect multiple times - it will reuse the existing connection
// if one exists
func Connect() (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	if connection == nil {
		var err error
		url := os.Getenv("DB_URL")
		authToken := os.Getenv("DB_TOKEN")
		if authToken != "" {
			// Handle remote connection to a Turso database
			url = fmt.Sprintf("%s?authToken=%s", url, authToken)
		}
		connection, err = sql.Open("libsql", url)
		if err != nil {
			return nil, err
		}

		// Verify connection by pinging
		if err := connection.Ping(); err != nil {
			connection.Close()
			connection = nil
			return nil, err
		}
	}

	return connection, nil
}

// Close closes the database connection if it's open
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if connection != nil {
		err := connection.Close()
		connection = nil
		return err
	}
	return nil
}
