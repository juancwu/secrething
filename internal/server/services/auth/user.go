package auth

import (
	"context"
	"time"

	"github.com/juancwu/konbini/internal/server/db"
	"github.com/juancwu/konbini/internal/server/utils"
)

func CreateUser(ctx context.Context, email, password string, name *string) (*db.User, error) {
	exists, err := ExistsUser(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, NewUserAlreadyExistsErr(email)
	}

	q, err := db.Query()
	if err != nil {
		return nil, err
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	timestamp := utils.FormatRFC3339NanoFixed(time.Now())

	user, err := q.NewUser(ctx, db.NewUserParams{
		Email:        email,
		PasswordHash: hashed,
		Name:         name,
		CreatedAt:    timestamp,
		UpdatedAt:    timestamp,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func ExistsUser(ctx context.Context, email string) (bool, error) {
	q, err := db.Query()
	if err != nil {
		// Default return true to avoid mistakenly proceed with other operations on error
		return true, err
	}
	_, err = q.ExistsUser(ctx, email)
	if db.IsNoRows(err) {
		return false, err
	}

	return true, nil
}
