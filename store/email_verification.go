package store

import (
	"database/sql"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	EMAIL_VERIFICATION_CODE_ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	EMAIL_VERIFICATION_CODE_LENGTH   = 20
)

// EmailVerification represents an email verification record in the database.
type EmailVerification struct {
	Id        int64
	Code      string
	UserId    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// Delete removes the email verification from the database.
// You must call tx.Commit for the deletion to take effect.
func (ev *EmailVerification) Delete(tx *sql.Tx) (sql.Result, error) {
	return tx.Exec("Delete FROM email_verifications WHERE id = $1;", ev.Id)
}

// Update is not implemented, but its defined to satisfy the Model interface.
func (ev *EmailVerification) Update(tx *sql.Tx) (sql.Result, error) {
	return nil, nil
}

// NewEmailVerification creates a new email verification record.
func NewEmailVerification(userId string) (*EmailVerification, error) {
	code, err := gonanoid.Generate(EMAIL_VERIFICATION_CODE_ALPHABET, EMAIL_VERIFICATION_CODE_LENGTH)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(time.Minute)
	row := db.QueryRow("INSERT INTO email_verifications (code, user_id, expires_at) VALUES ($1, $2, $3) RETURNING id, created_at;", code, userId, expiresAt)
	err = row.Err()
	if err != nil {
		return nil, err
	}
	ev := EmailVerification{
		Code:      code,
		UserId:    userId,
		ExpiresAt: expiresAt,
	}
	err = row.Scan(&ev.Id, &ev.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// GetEmailVerificationWithCode tries to retrieves an email verification record from the database with the given code.
func GetEmailVerificationWithCode(code string) (*EmailVerification, error) {
	row := db.QueryRow("SELECT id, code, user_id, expires_at, created_at FROM email_verifications WHERE code = $1;", code)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	ev := EmailVerification{}
	err = row.Scan(
		&ev.Id,
		&ev.Code,
		&ev.UserId,
		&ev.ExpiresAt,
		&ev.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// DeleteEmailVerificationWithUserId tries to removes any email verification code that has the given user id.
// IMPORTANT: this method does not use a transaction so deletion are UNSAFE.
func DeleteEmailVerificationWithUserId(uid string) (sql.Result, error) {
	return db.Exec("DELETE FROM email_verifications WHERE user_id = $1;", uid)
}
