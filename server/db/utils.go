package db

import (
	"database/sql"

	"github.com/rs/zerolog"
)

func RollabackWithLog(tx *sql.Tx, log *zerolog.Logger) {
	err := tx.Rollback()
	if err != nil {
		log.Error().Err(err).Msg("Failed to rollback.")
	}
}
