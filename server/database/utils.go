package database

import (
	"database/sql"

	"github.com/charmbracelet/log"
)

func Rollback(tx *sql.Tx, name string) {
	err := tx.Rollback()
	if err != nil {
		log.Errorf("Error rolling back changes (%s): %v\n", name, err)
	}
}
