package service

import (
	"database/sql"

	"github.com/matoous/go-nanoid/v2"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/env"
	usermodel "github.com/juancwu/konbini/server/models/user"
	"github.com/juancwu/konbini/server/utils"
)

func GetUserWithEmail(email string) (*usermodel.User, error) {
	user, err := usermodel.GetByEmail(email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.Logger().Info("No user found with email", "email", email)
			return nil, nil
		}
		utils.Logger().Errorf("Error getting user with email: %s, cause: %s\n", email, err)
		return nil, err
	}

	return user, nil
}

func RegisterUser(firstName, lastName, email, password string, tx *sql.Tx) (string, error) {
	utils.Logger().Info("Registering user with email", email)

	row := tx.QueryRow(
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, crypt($4, gen_salt($5))) RETURNING id;",
		firstName, lastName, email, password, env.Values().PASS_ENCRYPT_ALGO)
	if row.Err() != nil {
		utils.Logger().Errorf("Error resgitering user: %v\n", row.Err())
		return "", row.Err()
	}

	var id string
	err := row.Scan(&id)
	if err != nil {
		utils.Logger().Errorf("Error getting returning user id after insert: %v\n", err)
		return "", err
	}

	return id, nil
}

func CreateResetPasswordRecord(user *usermodel.User) (string, error) {
	// generate random id for reset password link
	linkId, err := gonanoid.New(32)
	if err != nil {
		utils.Logger().Errorf("Failed to generate random id for reset password link: %v\n", err)
		return "", err
	}
	utils.Logger().Infof("Random reset password link id: %s\n", linkId)

	// store reset password entry to generate a new reset password link that gets sent in the reset password email
	_, err = database.DB().Exec("INSERT INTO users_passwords_reset (user_id, link_id) VALUES ($1, $2);", user.Id, linkId)
	if err != nil {
		utils.Logger().Errorf("Failed to insert password reset link id into db: %v\n", err)
		return "", err
	}

	// send the email to reset the password

	return linkId, nil
}
