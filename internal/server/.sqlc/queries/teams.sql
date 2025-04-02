-- name: CreateTeam :one
INSERT INTO teams (
  team_id,
  name,
  description,
  created_by_user_id,
  created_at,
  updated_at
) VALUES (
  ?1, ?2, ?3, ?4, ?5, ?6
)
RETURNING team_id, name, description, created_by_user_id, created_at, updated_at;

-- name: GetTeamByID :one
SELECT team_id, name, description, created_by_user_id, created_at, updated_at
FROM teams
WHERE team_id = ?1;

-- name: GetTeamsByUserID :many
SELECT t.team_id, t.name, t.description, t.created_by_user_id, t.created_at, t.updated_at
FROM teams t
JOIN users_teams ut ON t.team_id = ut.team_id
WHERE ut.user_id = ?1;

-- name: UpdateTeam :one
UPDATE teams
SET name = ?2,
    description = ?3,
    updated_at = ?4
WHERE team_id = ?1
RETURNING team_id, name, description, created_by_user_id, created_at, updated_at;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE team_id = ?1;