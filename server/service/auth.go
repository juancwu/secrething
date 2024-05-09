package service

import (
	"database/sql"

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

func RegisterUser(firstName, lastName, email, password string, tx *sql.Tx) (int64, error) {
	utils.Logger().Info("Registering user with email", email)

	row := tx.QueryRow(
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, crypt($4, gen_salt($5))) RETURNING id;",
		firstName, lastName, email, password, env.Values().PASS_ENCRYPT_ALGO)
	if row.Err() != nil {
		utils.Logger().Errorf("Error resgitering user: %v\n", row.Err())
		return 0, row.Err()
	}

	var id int64
	err := row.Scan(&id)
	if err != nil {
		utils.Logger().Errorf("Error getting returning user id after insert: %v\n", err)
		return 0, err
	}

	return id, nil
}
