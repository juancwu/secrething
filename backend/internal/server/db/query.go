package db

import "database/sql"

func Query() (q *Queries, err error) {
	var conn *sql.DB
	conn, err = Connect()
	if err != nil {
		return
	}

	q = New(conn)

	return
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
