package db

import (
	"context"
	"database/sql"
	"fmt"
	"konbini/server/config"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
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
func (c *DBConnector) Connect() (*Connection, error) {
	cfg, err := config.Global()
	if err != nil {
		return nil, err
	}

	// connect based on the app environment
	var db *sql.DB
	if cfg.IsTesting() {
		db, err = connectToLocal(c.url)
	} else {
		db, err = connectToRemote(c.url)
	}
	if err != nil {
		return nil, err
	}

	return &Connection{conn: db}, nil
}

// The Connection struct implements or at least satisfies the sql.DB interface
// so that it can be used along with sqlc generated code and just normal usage
// of the sql.DB instance.
// This implementation is needed for properly handling local and remote connections
// without having to modify the code in handlers.
type Connection struct {
	conn *sql.DB
}

func (c *Connection) Begin() (*sql.Tx, error) {
	return c.conn.Begin()
}

func (c *Connection) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.conn.BeginTx(ctx, opts)
}

func (c *Connection) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.conn.ExecContext(ctx, query, args...)
}

func (c *Connection) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return c.conn.PrepareContext(ctx, query)
}

func (c *Connection) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.conn.QueryContext(ctx, query, args...)
}

func (c *Connection) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return c.conn.QueryRowContext(ctx, query, args...)
}

func (c *Connection) Close() error {
	cfg, err := config.Global()
	if err != nil {
		return err
	}
	if cfg.IsTesting() {
		return nil
	}
	return c.conn.Close()
}
