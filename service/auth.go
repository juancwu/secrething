package service

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/juancwu/konbini/database"
	"github.com/juancwu/konbini/model"
)

func GetUserByEmail(email string) (*model.UserModel, error) {
	log.Info("Getting user by email", email)
	row := database.DB().QueryRowContext(context.Background(), "SELECT id, first_name, last_name, email, pem_public_key FROM users WHERE email = ?;", email)
	if row.Err() != nil {
		log.Errorf("Error getting user with email: %s, cause: %s\n", email, row.Err())
		return nil, row.Err()
	}

	user := model.UserModel{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PemPublicKey)
	if err != nil {
		log.Errorf("Error scanning values for user: %v, function: GetUserByEmail\n", err)
		return nil, err
	}

	return &user, nil
}

func RegisterUser(firstName, lastName *string, email, pemPublicKey string) error {
	log.Info("Registering user with email", email)

	result, err := database.DB().ExecContext(context.Background(), "INSERT INTO users (first_name, last_name, email, pem_public_key) VALUES (?, ?, ?, ?);", firstName, lastName, email, pemPublicKey)
	if err != nil {
		log.Errorf("Error resgitering user: %v\n", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Errorf("Error getting rows inserted after registering user: %v\n", err)
	} else {
		log.Infof("Resgitered %d user(s)\n", count)
	}

	return nil
}
