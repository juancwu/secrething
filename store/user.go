package store

import (
	"os"
	"time"
)

// User represents a real user store in the database.
type User struct {
	Id            string
	Email         string
	Password      string
	Name          string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewUser creates a new user with the given information.
// This function will save the new user in the database.
func NewUser(email, password, name string) (*User, error) {
	user := User{
		Email:    email,
		Password: password,
		Name:     name,
	}
	row := db.QueryRow("INSERT INTO users (email, password, name) VALUES ($1, crypt($2, gen_salt($3)), $4) RETURNING id, email_verified, created_at, updated_at;", email, password, os.Getenv("PASS_ENCRYPT_ALGO"), name)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	err = row.Scan(
		&user.Id,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsUserWithEmail checks if there is a user with the given email stored in the database.
func ExistsUserWithEmail(email string) (bool, error) {
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email)
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
