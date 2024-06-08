package store

import (
	"time"
)

const (
	// EmailVerification Statuses
	STATUS_PENDING  = "PENDING"
	STATUS_SENT     = "SENT"
	STATUS_OPENED   = "OPENED"
	STATUS_VERIFIED = "VERIFIED"
	STATUS_FAILED   = "PENDING"

	EMAIL_VERIFICATION_CODE_LEN = 16
)

// EmailVerification represents how an email verification is stored and what can be accessed.
// This is used to manage emails that has been sent to for verifying the email of an account.
// Hence, there is also a expiration time as ExpiresAt.
type EmailVerification struct {
	Id            int64
	Code          string
	UserId        string
	ResendEmailId *string
	Status        string
	EmailSentAt   time.Time
	ExpiresAt     time.Time
	VerifiedAt    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CreateEmailVerification stores a new EmailVerification in the database.
func CreateEmailVerification(code, userId, resendEmailId string) error {
	// create expire time
	expiresAt := time.Now().Add(time.Minute * 30)

	// create new email verification
	row := db.QueryRow(
		"INSERT INTO email_verifications (code, user_id, resend_email_id, status, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		code,
		userId,
		resendEmailId,
		STATUS_PENDING,
		expiresAt,
	)
	err := row.Err()
	if err != nil {
		return err
	}

	// get returning id
	var id string
	err = row.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
