package service

import (
	"time"

	"github.com/charmbracelet/log"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/env"
)

type User struct {
	Id            int64
	FirstName     string
	LastName      string
	Password      string // password is always encrypted when fetched from db
	Email         string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func GetUserWithEmail(email string) (*User, error) {
	log.Info("Getting user by email", "email", email)
	user := User{}
	err := database.DB().
		QueryRow(
			"SELECT id, first_name, last_name, email, email_verified, created_at, updated_at FROM users WHERE email = $1;",
			email).
		Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.EmailVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Info("No user found with email", "email", email)
			return nil, nil
		}
		log.Errorf("Error getting user with email: %s, cause: %s\n", email, err)
		return nil, err
	}

	return &user, nil
}

func RegisterUser(firstName, lastName, email, password string) (int64, error) {
	log.Info("Registering user with email", email)

	row := database.DB().QueryRow(
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, crypt($4, gen_salt($5))) RETURNING id;",
		firstName, lastName, email, password, env.Values().PASS_ENCRYPT_ALGO)
	if row.Err() != nil {
		log.Errorf("Error resgitering user: %v\n", row.Err())
		return 0, row.Err()
	}

	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Errorf("Error getting returning user id after insert: %v\n", err)
		return 0, err
	}

	return id, nil
}
