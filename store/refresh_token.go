package store

import (
	"context"
	"database/sql"
	"time"
)

type RefreshToken struct {
	ID        string
	UserId    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

const new_REFRESH_TOKEN_SQL = `INSERT INTO refresh_tokens (user_id, token) VALUES ($1, $2) RETURNING id;`

// Creates a new refresh token in the database
func NewRefreshToken(ctx context.Context, db *sql.DB, userID, token string, expiresAt time.Time) (string, error) {
	var id string
	row := db.QueryRowContext(ctx, new_REFRESH_TOKEN_SQL, userID, token)
	err := row.Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

const get_REFRESH_TOKEN_BY_USERID_TOKEN = `SELECT
    id,
    user_id,
    token,
    expires_at,
    created_at
FROM refresh_tokens WHERE user_id = $1 AND token = $2;`

// Gets a refresh token record from the database by user id and token string.
func GetRefreshTokenByUserIDToken(ctx context.Context, db *sql.DB, userID, token string) (*RefreshToken, error) {
	rt := &RefreshToken{}
	row := db.QueryRowContext(ctx, get_REFRESH_TOKEN_BY_USERID_TOKEN, userID, token)
	err := row.Scan(
		&rt.ID,
		&rt.UserId,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return rt, nil
}
