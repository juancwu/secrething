package store

import (
	"database/sql"
	_ "embed"
	"time"
)

const (
	O_WRITE       int16 = 0b0000_0001
	O_SHARE       int16 = 0b0000_0010
	O_GRANT_SHARE int16 = 0b0000_0100
	O_DELETE      int16 = 0b0000_1000
)

//go:embed raw_sql/new_bento_permission.sql
var new_bento_permission_sql string

//go:embed raw_sql/get_bento_permission_user_bento_id.sql
var get_bento_permission_user_bento_id_sql string

// BentoPermission represents an entry for a bento permission in the database.
// It has a set of methods that help with the retrieval and manipulation of the bento permission.
type BentoPermission struct {
	Id          int64
	UserId      string
	BentoId     string
	Permissions int16
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// This will create a new bento permission in the database using the transaction provided.
//
// IMPORTANT: You must call tx.Commit() afterwards for changes to take effect.
func NewBentoPermissionTx(tx *sql.Tx, userId, bentoId string, permissions int16) (*BentoPermission, error) {
	perms := BentoPermission{
		UserId:      userId,
		BentoId:     bentoId,
		Permissions: permissions,
	}

	row := tx.QueryRow(new_bento_permission_sql, userId, bentoId, permissions)
	if err := row.Err(); err != nil {
		return nil, err
	}

	if err := row.Scan(
		&perms.Id,
		&perms.CreatedAt,
		&perms.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &perms, nil
}

// GetBentoPermissionByUserBentoId gets a bento permission that matches the given
// user and bento id.
func GetBentoPermissionByUserBentoId(userId, bentoId string) (*BentoPermission, error) {
	perms := new(BentoPermission)

	row := db.QueryRow(get_bento_permission_user_bento_id_sql, userId, bentoId)
	if err := row.Err(); err != nil {
		return nil, err
	}

	if err := row.Scan(
		&perms.Id,
		&perms.UserId,
		&perms.BentoId,
		&perms.Permissions,
		&perms.CreatedAt,
		&perms.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return perms, nil
}
