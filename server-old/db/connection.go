package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/tursodatabase/go-libsql"
)

// TursoConnector implements the Connector interface for Turso database
type TursoConnector struct {
	dbURL      string
	authToken  string
	connection *sql.DB
	mu         sync.RWMutex
}

// NewTursoConnector creates a new instance of TursoConnector
func NewTursoConnector(dbURL, authToken string) *TursoConnector {
	return &TursoConnector{
		dbURL:     dbURL,
		authToken: authToken,
	}
}

// Connect establishes a connection to the Turso database
// It's safe to call Connect multiple times - it will reuse the existing connection
// if one exists
func (c *TursoConnector) Connect() (*sql.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connection != nil {
		return c.connection, nil
	}

	// Open the database using the connector
	var err error
	var url string
	if c.authToken != "" {
		// Handle remote connection to a Turso database
		url = fmt.Sprintf("%s?authToken=%s", c.dbURL, c.authToken)
	} else {
		// Handle local connection to a SQLite database
		url = c.dbURL
	}
	c.connection, err = sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}

	// Verify connection by pinging
	if err := c.connection.Ping(); err != nil {
		c.connection.Close()
		c.connection = nil
		return nil, err
	}

	return c.connection, nil
}

// GetDB returns the current database connection
// Returns nil if no connection has been established
func (c *TursoConnector) GetDB() *sql.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connection
}

// Close closes the database connection if it's open
func (c *TursoConnector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connection != nil {
		err := c.connection.Close()
		c.connection = nil
		return err
	}
	return nil
}

// Transaction represents a database transaction
type Transaction struct {
	Tx *sql.Tx
}

// WithTransaction executes the given function within a transaction
// The transaction is committed if the function returns nil, otherwise it's rolled back
func WithTransaction(db *sql.DB, fn func(*Transaction) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Rollback is a no-op if the transaction is already committed
	defer tx.Rollback()

	if err := fn(&Transaction{Tx: tx}); err != nil {
		return err
	}

	return tx.Commit()
}
