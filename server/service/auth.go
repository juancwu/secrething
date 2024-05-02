package service

import (
	"time"

	"github.com/charmbracelet/log"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/env"
	"github.com/juancwu/konbini/server/sql"
)

type User struct {
	Id            int64
	FirstName     string
	LastName      string
	Password      string // password is always encrypted when fetched from db
	Email         string
	EmailVerified bool
	PubKey        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func GetUserWithEmail(email string) (*User, error) {
	log.Info("Getting user by email", "email", email)
	user := User{}
	err := database.DB().
		QueryRow(sql.GetUserWithEmail, email, env.Values().PGP_SYM_KEY).
		Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.EmailVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.PubKey,
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

	row := database.DB().QueryRow(sql.CreateUser, firstName, lastName, email, password)
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
