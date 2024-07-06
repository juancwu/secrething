// This file contains all things related to how the database store users
package store

import (
	"database/sql"
	"os"
	"time"
)

// User is someone who uses konbini services. This structure is a representation of
// a complete user fetched from the database.
type User struct {
	Id            string
	FirstName     string
	LastName      string
	Email         string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// UserExists checks if a user with the given email exists in the database or not.
func UserExists(email string) (bool, error) {
	var exists bool
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email)
	err := row.Err()
	if err != nil {
		return false, err
	}
	err = row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateUser inserts a new user into the database and returns the newly created user's id.
func CreateUser(email, password, firstName, lastName string) (string, error) {
	row := db.QueryRow(
		"INSERT INTO users (email, first_name, last_name, password) VALUES ($1, $2, $3, crypt($4, gen_salt($5))) RETURNING id;",
		email,
		firstName,
		lastName,
		password,
		os.Getenv("PASS_ENCRYPT_ALGO"),
	)

	err := row.Err()
	if err != nil {
		return "", err
	}

	// get new user id
	var id string
	err = row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetUserWithPasswordValidation tries to match a user with the given email and password.
// Ideal use of this function is for logging in a user with email and password.
// NOTE: This function will only return no if password matches for account with email and the email has been verified.
func GetUserWithPasswordValidation(email, password string) (*User, error) {
	row := db.QueryRow(
		`
        SELECT
            id,
            first_name,
            last_name,
            email,
            email_verified,
            created_at,
            updated_at
        FROM users
        WHERE email = $1 AND password = crypt($2, password);
        `,
		email,
		password,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	user := User{}
	err = row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserWithEmail retrieves a user with the given email.
// It will return an error if no row was found.
func GetUserWithEmail(email string) (*User, error) {
	row := db.QueryRow("SELECT id, email, first_name, last_name, email_verified, created_at, updated_at FROM users WHERE email = $1;", email)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	user := User{}
	err = row.Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserWithId retrieves a user with the given id.
// It will return an error if no row was found.
func GetUserWithId(id string) (*User, error) {
	row := db.QueryRow("SELECT id, email, first_name, last_name, email_verified, created_at, updated_at FROM users WHERE id = $1;", id)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	user := User{}
	err = row.Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// SetUserEmailVerifiedStatus updates the email_verified column of a user in the database.
func SetUserEmailVerifiedStatus(tx *sql.Tx, id string, status bool) error {
	_, err := tx.Exec(
		"UPDATE users SET email_verified = $1 WHERE id = $2;",
		status,
		id,
	)
	return err
}

// UpdateUserPasswordWithIdTx uses the given transaction to update the password of a user with a given id.
func UpdateUserPasswordWithIdTx(tx *sql.Tx, id, password string) (sql.Result, error) {
	return db.Exec("UPDATE users SET password = crypt($1, gen_salt($2)) WHERE id = $3;", password, os.Getenv("PASS_ENCRYPT_ALGO"), id)
}
