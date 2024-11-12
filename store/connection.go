package store

import (
	// builtin modules
	"context"
	"database/sql"
	"errors"
	"time"

	// package modules
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
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

// Close closes the database connection. Not intended to be used in production, but makes testing easier.
func Close() error {
	return db.Close()
}

// Ping makes sure that connection is still alives. It context.Background and timeouts in 5 seconds.
func Ping() error {
	if db == nil {
		return errors.New("db instance is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// create a channel to listen to when the ping finishes
	ch := make(chan error)

	go func() {
		err := db.PingContext(ctx)
		ch <- err
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-ch:
			return err
		}
	}
}

// StartTx begins a new transaction
func StartTx() (*sql.Tx, error) {
	return db.Begin()
}

// Rollback is a helper function to rollback if tx.Commit() returns an error
// and will log if rollback results in an error.
func Rollback(tx *sql.Tx, requestId string) error {
	if err := tx.Rollback(); err != nil {
		log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback")
		return err
	}
	return nil
}
