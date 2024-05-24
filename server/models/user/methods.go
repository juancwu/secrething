package usermodel

import (
	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/utils"
)

func GetByEmail(email string) (*User, error) {
	utils.Logger().Info("Getting user by email", "email", email)
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
		return nil, err
	}

	return &user, nil
}

func GetByEmailWithPassword(email, password string) (*User, error) {
	utils.Logger().Info("Getting user by email", "email", email)
	user := User{}
	err := database.DB().
		QueryRow(
			"SELECT id, first_name, last_name, email, email_verified, created_at, updated_at FROM users WHERE email = $1 AND password = crypt($2, password);",
			email,
			password,
		).
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
		return nil, err
	}

	return &user, nil
}

func GetById(id string) (*User, error) {
	utils.Logger().Info("Gettin user by id", "id", id)
	user := User{}
	err := database.DB().
		QueryRow(
			"SELECT id, first_name, last_name, email, email_verified, created_at, updated_at FROM users WHERE id = $1;",
			id,
		).
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
		return nil, err
	}

	return &user, nil
}

func IsRealUser(id string) (bool, error) {
	return utils.RowExists(database.DB(), "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", id)
}
