// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: groups.sql

package db

import (
	"context"
)

const addUserToGroup = `-- name: AddUserToGroup :exec
INSERT INTO users_groups (user_id, group_id, created_at) VALUES (?, ?, ?)
`

type AddUserToGroupParams struct {
	UserID    string `db:"user_id"`
	GroupID   string `db:"group_id"`
	CreatedAt string `db:"created_at"`
}

func (q *Queries) AddUserToGroup(ctx context.Context, arg AddUserToGroupParams) error {
	_, err := q.db.ExecContext(ctx, addUserToGroup, arg.UserID, arg.GroupID, arg.CreatedAt)
	return err
}

const existsGroupOwnedByUser = `-- name: ExistsGroupOwnedByUser :one
SELECT EXISTS(SELECT 1 FROM groups WHERE name = ? AND owner_id = ?)
`

type ExistsGroupOwnedByUserParams struct {
	Name    string `db:"name"`
	OwnerID string `db:"owner_id"`
}

func (q *Queries) ExistsGroupOwnedByUser(ctx context.Context, arg ExistsGroupOwnedByUserParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, existsGroupOwnedByUser, arg.Name, arg.OwnerID)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const existsGroupWithIdOwnedByUser = `-- name: ExistsGroupWithIdOwnedByUser :one
SELECT EXISTS(SELECT 1 FROM groups WHERE id = ? AND owner_id = ?)
`

type ExistsGroupWithIdOwnedByUserParams struct {
	ID      string `db:"id"`
	OwnerID string `db:"owner_id"`
}

func (q *Queries) ExistsGroupWithIdOwnedByUser(ctx context.Context, arg ExistsGroupWithIdOwnedByUserParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, existsGroupWithIdOwnedByUser, arg.ID, arg.OwnerID)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getGroupByIDOwendByUser = `-- name: GetGroupByIDOwendByUser :one
SELECT id, name, owner_id, created_at, updated_at FROM groups WHERE id = ? AND owner_id = ?
`

type GetGroupByIDOwendByUserParams struct {
	ID      string `db:"id"`
	OwnerID string `db:"owner_id"`
}

func (q *Queries) GetGroupByIDOwendByUser(ctx context.Context, arg GetGroupByIDOwendByUserParams) (Group, error) {
	row := q.db.QueryRowContext(ctx, getGroupByIDOwendByUser, arg.ID, arg.OwnerID)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OwnerID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const newGroup = `-- name: NewGroup :one
INSERT INTO groups (name, owner_id, created_at, updated_at)
VALUES (?, ?, ?, ?)
RETURNING id
`

type NewGroupParams struct {
	Name      string `db:"name"`
	OwnerID   string `db:"owner_id"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (q *Queries) NewGroup(ctx context.Context, arg NewGroupParams) (string, error) {
	row := q.db.QueryRowContext(ctx, newGroup,
		arg.Name,
		arg.OwnerID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const newGroupInvitation = `-- name: NewGroupInvitation :one
INSERT INTO group_invitations
(user_id, group_id, created_at, expires_at)
VALUES (?, ?, ?, ?)
RETURNING id
`

type NewGroupInvitationParams struct {
	UserID    string `db:"user_id"`
	GroupID   string `db:"group_id"`
	CreatedAt string `db:"created_at"`
	ExpiresAt string `db:"expires_at"`
}

func (q *Queries) NewGroupInvitation(ctx context.Context, arg NewGroupInvitationParams) (string, error) {
	row := q.db.QueryRowContext(ctx, newGroupInvitation,
		arg.UserID,
		arg.GroupID,
		arg.CreatedAt,
		arg.ExpiresAt,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const removeGroupByID = `-- name: RemoveGroupByID :exec
DELETE FROM groups WHERE id = ?
`

func (q *Queries) RemoveGroupByID(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, removeGroupByID, id)
	return err
}

const removeUserFromGroup = `-- name: RemoveUserFromGroup :exec
DELETE FROM users_groups WHERE user_id = ? AND group_id = ?
`

type RemoveUserFromGroupParams struct {
	UserID  string `db:"user_id"`
	GroupID string `db:"group_id"`
}

func (q *Queries) RemoveUserFromGroup(ctx context.Context, arg RemoveUserFromGroupParams) error {
	_, err := q.db.ExecContext(ctx, removeUserFromGroup, arg.UserID, arg.GroupID)
	return err
}
