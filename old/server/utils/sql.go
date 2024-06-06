package utils

import (
	"database/sql"
)

func RowExists(db *sql.DB, query string, args ...interface{}) (bool, error) {
	var exists bool
	row := db.QueryRow(query, args...)
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}
