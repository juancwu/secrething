package store

import "database/sql"

type Model interface {
	// DELETE is used to delete the row in the database.
	DELETE(tx *sql.Tx) error
	// UPDATE will save the current model values into the database.
	UPDATE(tx *sql.Tx) (sql.Result, error)
}
