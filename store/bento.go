package store

import (
	"os"
	"time"

	"go.uber.org/zap"
)

// Bento represents a row in the bentos table.
type Bento struct {
	Id      string
	Name    string
	OwnerId string
	PubKey  []byte
	// IsShared is not an official column in the database but it should be filled in when querying a bento
	IsShared  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	// Description on the permissions with their corresponding integer value.
	// Note that permissions aggregate to as you go down the permissions list.
	// How to read the permissions:
	// A permission integer value is of at most 4 numeric places -> xxxx
	// where each place represents share,delete,write,read respectively.
	// At any given place, if the value is 7, then its given ALL permissions of said place.
	// Giving write:all (0077) permissions to a user also gives read:all (0007) permissions.
	// Now the why use integer values instead of a more readable string version.
	// Reason: saves space in the database if we store a SMALLINT that is of 2 bytes.
	// read:all 0007 - read all
	BENTO_READ_ALL = 7
	// write:all 0077 - write all
	BENTO_WRITE_ALL = 77
	// write:name 0017 - change bento name
	BENTO_WRITE_NAME = 17
	// write:ingridient:key 0027 - change bento ingridient key
	BENTO_WRITE_INGRIDIENT_KEY = 27
	// delete:all 0777 - can delete the entire bento
	BENTO_DELETE_ALL = 777
	// delete:ingridients 0177 - can only delete ingridients but not the entire bento
	BENTO_DELETE_INGRIDIENTS = 177
	// admin:all 7777 - can share everything and super admin access to the bento
	BENTO_ADMIN_ALL = 7777
	// share:read:all 1777 - can share bento with read access
	BENTO_SHARE_READ_ALL = 1777
	// share:write:all 2777 - can share with read/write access
	BENTO_SHARE_WRITE_ALL = 2777
	// share:write:name 3777 - can share with read/write but only to bento's name access
	BENTO_SHARE_WRITE_NAME = 3777
	// share:write:ingridient:key 4777 - can share with write access to individual ingridient key name
	BENTO_SHARE_WRITE_INGRIDIENT_KEY = 4777
	// share:delete:all 5777 - can share with delete access
	BENTO_SHARE_DELETE_ALL = 5777
	// share:delete:ingridients 6777 - can share with delete ingridients access
	BENTO_SHARE_DELETE_INGRIDIENTS = 6777
)

// PrepBento creates a new bento in the bentos table but does not create new entries (ingridients).
func PrepBento(name, ownerId, pubKey string) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	var bentoId string
	row := tx.QueryRow(
		"INSERT INTO bentos (name, owner_id, pub_key) VALUES ($1, $2, pgp_sym_encrypt($3, $4)) RETURNING id;",
		name,
		ownerId,
		pubKey,
		os.Getenv("PGP_SYM_KEY"),
	)
	err = row.Err()
	if err != nil {
		return "", err
	}
	err = row.Scan(&bentoId)
	if err != nil {
		return "", err
	}

	// should create a bento user permission entry to normalize the access to bento features
	_, err = tx.Exec(
		"INSERT INTO bento_user_permissions (user_id, bento_id, permissions) VALUES ($1, $2, $3);",
		ownerId,
		bentoId,
		BENTO_ADMIN_ALL,
	)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return bentoId, nil
}

// AddIngridient adds a new ingridient entry for a bento.
// This function does not check write permissions.
func AddIngridient(bentoId, key, value string) error {
	_, err := db.Exec(
		"INSERT INTO bento_ingridients (key, value, bento_id) VALUES ($1, $2, $3);",
		key,
		value,
		bentoId,
	)
	return err
}

// GetBento retrieves a bento with the given id. It will try to retrieve both owned and shared bentos.
func GetBento(bentoId string) (*Bento, error) {
	zap.L().Info("getting bento", zap.String("bento_id", bentoId))
	row := db.QueryRow(
		`
        SELECT
            id,
            name,
            owner_id,
            pub_key,
            created_at,
            updated_at,
            false as is_shared
        FROM bentos
        WHERE id = $1
        UNION
        SELECT
            b.id,
            b.name,
            b.owner_id,
            b.pub_key,
            b.created_at,
            b.updated_at,
            true as is_shared
        FROM bentos as b
        JOIN shared_bentos as sb ON sb.shared_bento_id = b.id
        WHERE sb.shared_bento_id = $1;
        `,
		bentoId,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	bento := Bento{}
	err = row.Scan(
		&bento.Id,
		&bento.Name,
		&bento.OwnerId,
		&bento.PubKey,
		&bento.CreatedAt,
		&bento.UpdatedAt,
		&bento.IsShared,
	)
	if err != nil {
		return nil, err
	}

	return &bento, nil
}
