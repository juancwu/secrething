// This file contains all things related to how the database store users
package store

import (
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
func CreateUser(email, password string, firstName, lastName *string) (string, error) {
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
