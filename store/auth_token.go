package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

// Checks if the token with the given tokenId exists and its valid
// in the database. This is because users can invalidate all access tokens
// for a given account and when resetting passwords, the tokens are invalidated too.
func IsTokenValid(tokenId, tokenType string) error {
	row := db.QueryRow("SELECT expires_at FROM auth_tokens WHERE id = $1 AND token_type = $2;", tokenId, tokenType)
	if err := row.Err(); err != nil {
		return err
	}
	var expiresAt time.Time
	if err := row.Scan(&expiresAt); err != nil {
		return err
	}
	if time.Now().After(expiresAt) {
		// spawn a go routine to not block the response to the client
		go func(id string) {
			if _, err := db.Exec("DELETE FROM auth_tokens WHERE id = $1;", id); err != nil {
				log.Error().Err(err).Str("token_id", id).Msg("Failed to delete expired token in auth_tokens.")
			}
		}(tokenId)
		return errors.New("Token exists in database but it has expired")
	}
	return nil
}

// Deletes all the tokens, access and refresh tokens owned by the given user with the userId.
func DeleteTokensOwnedByUser(tx *sql.Tx, userId string) (sql.Result, error) {
	return tx.Exec("DELETE FROM auth_tokens WHERE user_id = $1;", userId)
}
