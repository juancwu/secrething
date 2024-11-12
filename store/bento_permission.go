package store

import (
	"database/sql"
	_ "embed"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	O_NO_PERMS          int = 0b0000_0000_0000_0000
	O_WRITE             int = 0b0000_0000_0000_0001
	O_SHARE             int = 0b0000_0000_0000_0010
	O_GRANT_SHARE       int = 0b0000_0000_0000_0100
	O_DELETE            int = 0b0000_0000_0000_1000
	O_WRITE_INGRIDIENT  int = 0b0000_0000_0001_0000
	O_DELETE_INGRIDIENT int = 0b0000_0000_0010_0000
	O_RENAME_INGRIDIENT int = 0b0000_0000_0100_0000
	O_RENAME_BENTO      int = 0b0000_0000_1000_0000
	O_REVOKE_SHARE      int = 0b0000_0001_0000_0000
	O_OWNER             int = 0b1000_0000_0000_0000

	// text represnetation of the perms

	// This represents all the permissions a requesting user can grant to a target user.
	S_ALL               string = "all"
	S_WRITE             string = "write"
	S_DELETE            string = "delete"
	S_SHARE             string = "share"
	S_RENAME_BENTO      string = "rename_bento"
	S_RENAME_INGRIDIENT string = "rename_ingridient"
	S_WRITE_INGRIDIENT  string = "write_ingridient"
	S_DELETE_INGRIDIENT string = "delete_ingridient"
	S_REVOKE_SHARE      string = "revoke_share"
)

// An map to take the integer permission value using text
var TextToBinPerms map[string]int = map[string]int{
	S_WRITE:             O_WRITE,
	S_DELETE:            O_DELETE,
	S_SHARE:             O_SHARE,
	S_RENAME_BENTO:      O_RENAME_BENTO,
	S_RENAME_INGRIDIENT: O_RENAME_INGRIDIENT,
	S_WRITE_INGRIDIENT:  O_WRITE_INGRIDIENT,
	S_DELETE_INGRIDIENT: O_DELETE_INGRIDIENT,
	S_REVOKE_SHARE:      O_REVOKE_SHARE,
}

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
	Permissions int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// This will create a new bento permission in the database using the transaction provided.
//
// IMPORTANT: You must call tx.Commit() afterwards for changes to take effect.
func NewBentoPermissionTx(tx *sql.Tx, userId, bentoId string, permissions int) (*BentoPermission, error) {
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

// Checks if there already exists a bento permission for the given user and bento.
func ExistsBentoPermissionByUserBentoId(userId, bentoId string) (bool, error) {
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM bento_permissions WHERE user_id = $1 AND bento_id = $2)", userId, bentoId)
	err := row.Err()
	if err != nil {
		return false, err
	}
	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Deletes an entry in bento_permissions by matching the given id.
//
// IMPORTANT: THIS PROCESS IS NOT REVERSABLE!
func DeleteBentoPermissionById(permsId int64) error {
	res, err := db.Exec("DELETE FROM bento_permissions WHERE id = $1;", permsId)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get rows affected after performing DELETE on bento_permissions")
	} else if n > 1 {
		log.Warn().Int64("perms_id", permsId).Msg("More than 1 permissions entry has been deleted.")
	}
	return nil
}
