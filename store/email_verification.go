package store

import (
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
