package store

import (
	"os"
	"time"
)

// Bento represents a row in the bentos table.
type Bento struct {
	Id        string
	Name      string
	OwnerId   string
	PubKey    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PrepBento creates a new bento in the bentos table but does not create new entries (ingridients).
func PrepBento(name, ownerId, pubKey string) (string, error) {
	var id string
	row := db.QueryRow(
		"INSERT INTO bentos (name, owner_id, pub_key) VALUES ($1, $2, pgp_sym_encrypt($3, $4)) RETURNING id;",
		name,
		ownerId,
		pubKey,
		os.Getenv("PGP_SYM_KEY"),
	)
	err := row.Err()
	if err != nil {
		return "", err
	}

	err = row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
