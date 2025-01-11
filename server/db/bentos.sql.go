// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: bentos.sql

package db

import (
	"context"
)

const addIngredientToBento = `-- name: AddIngredientToBento :exec
INSERT INTO bento_ingredients (bento_id, name, value, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
`

type AddIngredientToBentoParams struct {
	BentoID   string `db:"bento_id"`
	Name      string `db:"name"`
	Value     []byte `db:"value"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (q *Queries) AddIngredientToBento(ctx context.Context, arg AddIngredientToBentoParams) error {
	_, err := q.db.ExecContext(ctx, addIngredientToBento,
		arg.BentoID,
		arg.Name,
		arg.Value,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const existsBentoWithNameOwnedByUser = `-- name: ExistsBentoWithNameOwnedByUser :one
SELECT EXISTS(SELECT 1 FROM bentos WHERE name = ? AND user_id = ?)
`

type ExistsBentoWithNameOwnedByUserParams struct {
	Name   string `db:"name"`
	UserID string `db:"user_id"`
}

func (q *Queries) ExistsBentoWithNameOwnedByUser(ctx context.Context, arg ExistsBentoWithNameOwnedByUserParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, existsBentoWithNameOwnedByUser, arg.Name, arg.UserID)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getBentoWithIDOwnedByUser = `-- name: GetBentoWithIDOwnedByUser :one
SELECT id, user_id, name, created_at, updated_at FROM bentos WHERE id = ? AND user_id = ?
`

type GetBentoWithIDOwnedByUserParams struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
}

func (q *Queries) GetBentoWithIDOwnedByUser(ctx context.Context, arg GetBentoWithIDOwnedByUserParams) (Bento, error) {
	row := q.db.QueryRowContext(ctx, getBentoWithIDOwnedByUser, arg.ID, arg.UserID)
	var i Bento
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const newBento = `-- name: NewBento :one
INSERT INTO bentos (user_id, name, created_at, updated_at)
VALUES (?, ?, ?, ?) RETURNING id
`

type NewBentoParams struct {
	UserID    string `db:"user_id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (q *Queries) NewBento(ctx context.Context, arg NewBentoParams) (string, error) {
	row := q.db.QueryRowContext(ctx, newBento,
		arg.UserID,
		arg.Name,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const setBentoIngredient = `-- name: SetBentoIngredient :exec
INSERT INTO bento_ingredients (bento_id, name, value, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT DO UPDATE SET
    value = excluded.value,
    updated_at = excluded.updated_at
`

type SetBentoIngredientParams struct {
	BentoID   string `db:"bento_id"`
	Name      string `db:"name"`
	Value     []byte `db:"value"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (q *Queries) SetBentoIngredient(ctx context.Context, arg SetBentoIngredientParams) error {
	_, err := q.db.ExecContext(ctx, setBentoIngredient,
		arg.BentoID,
		arg.Name,
		arg.Value,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}
