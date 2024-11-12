package store

import "database/sql"

type Model interface {
	// Delete is used to delete the row in the database.
	Delete(tx *sql.Tx) (sql.Result, error)
	// Update will save the current model values into the database.
	Update(tx *sql.Tx) (sql.Result, error)
}
