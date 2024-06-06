package database

import (
	"database/sql"

	"github.com/juancwu/konbini/server/utils"
)

func Rollback(tx *sql.Tx, name string) {
	err := tx.Rollback()
	if err != nil {
		utils.Logger().Errorf("Error rolling back changes (%s): %v\n", name, err)
	}
}
