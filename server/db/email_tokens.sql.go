// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: email_tokens.sql

package db

import (
	"context"
)

const createEmailToken = `-- name: CreateEmailToken :one
INSERT INTO email_tokens
(user_id, token_salt, created_at, expires_at)
VALUES (?, ?, ?, ?)
RETURNING id
`

type CreateEmailTokenParams struct {
	UserID    string
	TokenSalt []byte
	CreatedAt string
	ExpiresAt string
}

func (q *Queries) CreateEmailToken(ctx context.Context, arg CreateEmailTokenParams) (string, error) {
	row := q.db.QueryRowContext(ctx, createEmailToken,
		arg.UserID,
		arg.TokenSalt,
		arg.CreatedAt,
		arg.ExpiresAt,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const deleteAllEmailTokensByUserId = `-- name: DeleteAllEmailTokensByUserId :many
DELETE FROM email_tokens
WHERE user_id = ?
RETURNING id
`

func (q *Queries) DeleteAllEmailTokensByUserId(ctx context.Context, userID string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, deleteAllEmailTokensByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
