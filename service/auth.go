package service

import (
	"github.com/charmbracelet/log"

	"github.com/juancwu/konbini/database"
	"github.com/juancwu/konbini/model"
)

func GetUserByEmail(email string) (*model.UserModel, error) {
	log.Info("Getting user by email", "email", email)
	user := model.UserModel{}
	err := database.DB().QueryRow("SELECT id, first_name, last_name, email, pem_public_key FROM users WHERE email = $1", email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PemPublicKey)
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

func RegisterUser(firstName, lastName, email, pemPublicKey string) (int64, error) {
	log.Info("Registering user with email", email)

	row := database.DB().QueryRow("INSERT INTO users (first_name, last_name, email, pem_public_key) VALUES ($1, $2, $3, $4) RETURNING id;", firstName, lastName, email, pemPublicKey)
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
