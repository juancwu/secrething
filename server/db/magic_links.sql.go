// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: magic_links.sql

package db

import (
	"context"
)

const createMagicLink = `-- name: CreateMagicLink :one
INSERT INTO magic_links
(user_id, state, created_at, expires_at)
VALUES
(?, ?, ?, ?)
RETURNING id
`

type CreateMagicLinkParams struct {
	UserID    string
	State     string
	CreatedAt string
	ExpiresAt string
}

func (q *Queries) CreateMagicLink(ctx context.Context, arg CreateMagicLinkParams) (string, error) {
	row := q.db.QueryRowContext(ctx, createMagicLink,
		arg.UserID,
		arg.State,
		arg.CreatedAt,
		arg.ExpiresAt,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const getMagicLink = `-- name: GetMagicLink :one
SELECT id, user_id, state, created_at, expires_at
FROM magic_links
WHERE id = ? AND user_id = ? AND state = ?
`

type GetMagicLinkParams struct {
	ID     string
	UserID string
	State  string
}

func (q *Queries) GetMagicLink(ctx context.Context, arg GetMagicLinkParams) (MagicLink, error) {
	row := q.db.QueryRowContext(ctx, getMagicLink, arg.ID, arg.UserID, arg.State)
	var i MagicLink
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.State,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const removeMagicLink = `-- name: RemoveMagicLink :exec
DELETE FROM magic_links
WHERE id = ? AND user_id = ? AND state = ?
`

type RemoveMagicLinkParams struct {
	ID     string
	UserID string
	State  string
}

func (q *Queries) RemoveMagicLink(ctx context.Context, arg RemoveMagicLinkParams) error {
	_, err := q.db.ExecContext(ctx, removeMagicLink, arg.ID, arg.UserID, arg.State)
	return err
}
