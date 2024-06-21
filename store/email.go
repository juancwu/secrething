package store

import (
	"database/sql"
	"time"
)

const (
	EMAIL_VERIFICATION_CODE_LEN      = 16
	EMAIL_VERIFICATION_CODE_CHR_POOL = "abcdefghijklnmopqsrtvwxyzABCDEFGHIJKLNMOPQSRTVWXYZ1234567890"
)

// EmailVerification represents how an email verification is stored and what can be accessed.
// This is used to manage emails that has been sent to for verifying the email of an account.
// Hence, there is also a expiration time as ExpiresAt.
type EmailVerification struct {
	Id        int64
	Code      string
	UserId    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// CreateEmailVerification stores a new EmailVerification in the database.
func CreateEmailVerification(code, userId string) (int64, error) {
	// create expire time
	expiresAt := time.Now().Add(time.Minute * 30)

	// create new email verification
	row := db.QueryRow(
		"INSERT INTO email_verifications (code, user_id, expires_at) VALUES ($1, $2, $3) RETURNING id;",
		code,
		userId,
		expiresAt,
	)
	err := row.Err()
	if err != nil {
		return 0, err
	}

	// get returning id
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetEmailVerificationWithCode retrieves an email verification record from the database with the given code.
func GetEmailVerificationWithCode(code string) (*EmailVerification, error) {
	row := db.QueryRow(
		"SELECT id, code, user_id, expires_at, created_at FROM email_verifications WHERE code = $1;",
		code,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	emailVerification := EmailVerification{}
	err = row.Scan(
		&emailVerification.Id,
		&emailVerification.Code,
		&emailVerification.UserId,
		&emailVerification.ExpiresAt,
		&emailVerification.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &emailVerification, nil
}

func GetEmailVerificationWithUserId(uid string) (*EmailVerification, error) {
	row := db.QueryRow(
		"SELECT id, code, user_id, expires_at, created_at FROM email_verifications WHERE user_id = $1;",
		uid,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	emailVerification := EmailVerification{}
	err = row.Scan(
		&emailVerification.Id,
		&emailVerification.Code,
		&emailVerification.UserId,
		&emailVerification.ExpiresAt,
		&emailVerification.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &emailVerification, nil
}

// DeleteEmailVerificationTx uses a transaction to remove an email verification from the database.
// Pass the PK of the email verification that needs to be deleted.
func DeleteEmailVerificationTx(tx *sql.Tx, id int64) error {
	_, err := tx.Exec("DELETE FROM email_verifications WHERE id = $1;", id)
	return err
}

// DeleteEmailVerification removes an email verification from the database.
// Pass the PK of the email verification that needs to be deleted.
func DeleteEmailVerification(id int64) error {
	_, err := db.Exec("DELETE FROM email_verifications WHERE id = $1;", id)
	return err
}

// DeleteAllEmailVerificationFromUser deletes all the email verifications from the database that are
// linked to a user. Pass the user id.
func DeleteAllEmailVerificationFromUser(uid string) (int64, error) {
	res, err := db.Exec("DELETE FROM email_verifications WHERE user_id = $1;", uid)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
