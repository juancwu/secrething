package store

import (
	"database/sql"
	"time"
)

// PasswordReset represents the information needed to complete a password reset process.
type PasswordReset struct {
	Id        int64
	UserId    string
	ResetCode string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// SavePasswordResetCode stores a password reset code into the database.
func SavePasswordResetCode(resetCode, userId string, expiresAt time.Time) (int64, error) {
	row := db.QueryRow("INSERT INTO password_resets (reset_code, user_id, expires_at) VALUES ($1, $2, $3) RETURNING id;", resetCode, userId, expiresAt)
	err := row.Err()
	if err != nil {
		return 0, err
	}
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// ExistsPasswordResetForUser checks is there is already an existing password reset record in the
// database for the given user. Only one should exists at a time per user. It returns false if error.
func ExistsPasswordResetForUser(uid string) (bool, error) {
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM password_resets WHERE user_id = $1);", uid)
	err := row.Err()
	if err != nil {
		return false, err
	}
	var exists bool
	err = row.Scan(&exists)
	return exists, nil
}

// DeletePasswordResetByUserId removes all password reset records in the database with the given user id.
func DeletePasswordResetByUserId(uid string) error {
	_, err := db.Exec("DELETE FROM password_resets WHERE user_id = $1;", uid)
	return err
}

// GetPasswordResetForUser retrieves a PasswordReset struct.
// It will return an error if nothing was found.
func GetPasswordResetForUser(uid string) (*PasswordReset, error) {
	row := db.QueryRow("SELECT id, user_id, reset_code, expires_at, created_at FROM password_resets WHERE user_id = $%1;", uid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	pr := PasswordReset{}
	err = row.Scan(
		&pr.Id,
		&pr.UserId,
		&pr.ResetCode,
		&pr.ExpiresAt,
		&pr.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// DeletePasswordReset deletes a password reset record from the database that matches a given id.
func DeletePasswordReset(id int64) (sql.Result, error) {
	return db.Exec("DELETE FROM password_resets WHERE id = $1;", id)
}

// DeletePasswordResetTx uses the given transaction to delete a password reset record from the dabatase that matches a given id.
func DeletePasswordResetTx(tx *sql.Tx, id int64) (sql.Result, error) {
	return tx.Exec("DELETE FROM password_resets WHERE id = $1;", id)
}
