package store

import (
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
