/*
This is a collection of testing utilities for black box testing.
*/
package test

import (
	"context"
	"konbini/server/config"
	"konbini/server/db"
	"konbini/server/utils"
	"time"
)

type testUser struct {
	Id        string
	Email     string
	Password  string
	NickName  string
	TokenSalt []byte
}

var (
	testUserOne testUser = testUser{
		Email:    "userone@mail.com",
		Password: "useronepassword",
		NickName: "User One",
	}
	testUserTwo testUser = testUser{
		Email:    "usertwo@mail.com",
		Password: "usertwopassword",
	}
	testUsers []testUser = []testUser{
		testUserOne,
		testUserTwo,
	}
)

func seedWithTestUsers(cfg *config.Config) error {
	dbUrl, dbAuthToken := cfg.GetDatabaseConfig()
	connector := db.NewConnector(dbUrl, dbAuthToken)
	connection, err := connector.Connect()
	if err != nil {
		return err
	}
	defer connection.Close()
	queries := db.New(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, user := range testUsers {
		hashed, err := utils.GeneratePasswordHash(user.Password)
		if err != nil {
			return err
		}
		salt, err := utils.RandomBytes(16)
		if err != nil {
			return err
		}
		now := time.Now().UTC().Format(time.RFC3339)
		row, err := queries.CreateUser(ctx, db.CreateUserParams{
			Email:     user.Email,
			Password:  hashed,
			Nickname:  user.NickName,
			TokenSalt: salt,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return err
		}
		user.TokenSalt = salt
		user.Id = row.ID
	}

	return nil
}

func removeSeededTestUsers(cfg *config.Config) error {
	dbUrl, dbAuthToken := cfg.GetDatabaseConfig()
	connector := db.NewConnector(dbUrl, dbAuthToken)
	connection, err := connector.Connect()
	if err != nil {
		return err
	}
	defer connection.Close()
	queries := db.New(connection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, user := range testUsers {
		err := queries.DeleteUserById(ctx, user.Id)
		if err != nil {
			return err
		}
	}

	return nil
}
