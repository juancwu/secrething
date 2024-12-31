package services

import (
	"context"
	"konbini/server/db"
)

type RegisterUserParams struct {
	Email    string
	Password string
	NickName string
}

func RegisterUser(ctx context.Context, queries *db.Queries, params RegisterUserParams) (*db.User, error) {
	key, err := GetRandomJWTKey()
	if err != nil {
		return nil, err
	}

	createUserParams := db.CreateUserParams{
		Email: params.Email,
	}

	queries.CreateUser(ctx, db.CreateUserParams{})

	return nil, nil
}
