package store

import (
	"database/sql"
	_ "embed"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// Represents a row in the db for password reset codes.
type PasswordResetCode struct {
	Id        int64
	Code      string
	UserId    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Delete the password reset code from the database.
func (p *PasswordResetCode) Delete(tx *sql.Tx) (sql.Result, error) {
	return tx.Exec("DELETE FROM password_reset_codes WHERE id = $1;", p.Id)
}

//go:embed raw_sql/new_or_update_password_reset_code.sql
var newOrUpdatePasswordResetCodeSQL string

func NewOrUpdatePasswordResetCode(uid string) (*PasswordResetCode, error) {
	code, err := gonanoid.New(6)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Minute * 3)

	row := db.QueryRow(newOrUpdatePasswordResetCodeSQL, uid, code, expiresAt)
	if err := row.Err(); err != nil {
		return nil, err
	}

	prc := PasswordResetCode{
		Code:      code,
		UserId:    uid,
		ExpiresAt: expiresAt,
	}

	if err := row.Scan(&prc.Id, &prc.CreatedAt, &prc.UpdatedAt); err != nil {
		return nil, err
	}

	return &prc, nil
}

//go:embed raw_sql/get_password_reset_code_by_user_id.sql
var getPasswordResetCodeByUserIdSQL string

// Get a password reset code by using the given user id.
func GetPasswordResetCodeByUserId(uid string) (*PasswordResetCode, error) {
	row := db.QueryRow(getPasswordResetCodeByUserIdSQL, uid)
	if err := row.Err(); err != nil {
		return nil, err
	}
	prc := PasswordResetCode{}
	if err := row.Scan(&prc.Id, &prc.Code, &prc.UserId, &prc.ExpiresAt, &prc.CreatedAt, &prc.UpdatedAt); err != nil {
		return nil, err
	}
	return &prc, nil
}
