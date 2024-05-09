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
