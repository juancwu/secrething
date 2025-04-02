// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users_teams.sql

package db

import (
	"context"
)

const addUserToTeam = `-- name: AddUserToTeam :one
INSERT INTO users_teams (
  user_id,
  team_id,
  created_at,
  updated_at
) VALUES (
  ?1, ?2, ?3, ?4
)
RETURNING user_id, team_id, created_at, updated_at
`

type AddUserToTeamParams struct {
	UserID    string `db:"user_id" json:"user_id"`
	TeamID    string `db:"team_id" json:"team_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

func (q *Queries) AddUserToTeam(ctx context.Context, arg AddUserToTeamParams) (UsersTeam, error) {
	row := q.db.QueryRowContext(ctx, addUserToTeam,
		arg.UserID,
		arg.TeamID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i UsersTeam
	err := row.Scan(
		&i.UserID,
		&i.TeamID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTeamMembers = `-- name: GetTeamMembers :many
SELECT user_id, team_id, created_at, updated_at
FROM users_teams
WHERE team_id = ?1
`

func (q *Queries) GetTeamMembers(ctx context.Context, teamID string) ([]UsersTeam, error) {
	rows, err := q.db.QueryContext(ctx, getTeamMembers, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UsersTeam
	for rows.Next() {
		var i UsersTeam
		if err := rows.Scan(
			&i.UserID,
			&i.TeamID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const isUserInTeam = `-- name: IsUserInTeam :one
SELECT EXISTS (
  SELECT 1
  FROM users_teams
  WHERE user_id = ?1 AND team_id = ?2
)
`

type IsUserInTeamParams struct {
	UserID string `db:"user_id" json:"user_id"`
	TeamID string `db:"team_id" json:"team_id"`
}

func (q *Queries) IsUserInTeam(ctx context.Context, arg IsUserInTeamParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, isUserInTeam, arg.UserID, arg.TeamID)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const removeAllUsersFromTeam = `-- name: RemoveAllUsersFromTeam :exec
DELETE FROM users_teams
WHERE team_id = ?1
`

func (q *Queries) RemoveAllUsersFromTeam(ctx context.Context, teamID string) error {
	_, err := q.db.ExecContext(ctx, removeAllUsersFromTeam, teamID)
	return err
}

const removeUserFromTeam = `-- name: RemoveUserFromTeam :exec
DELETE FROM users_teams
WHERE user_id = ?1 AND team_id = ?2
`

type RemoveUserFromTeamParams struct {
	UserID string `db:"user_id" json:"user_id"`
	TeamID string `db:"team_id" json:"team_id"`
}

func (q *Queries) RemoveUserFromTeam(ctx context.Context, arg RemoveUserFromTeamParams) error {
	_, err := q.db.ExecContext(ctx, removeUserFromTeam, arg.UserID, arg.TeamID)
	return err
}
