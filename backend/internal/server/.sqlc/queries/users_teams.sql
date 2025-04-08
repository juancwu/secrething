-- name: AddUserToTeam :one
INSERT INTO users_teams (
  user_id,
  team_id,
  created_at,
  updated_at
) VALUES (
  ?1, ?2, ?3, ?4
)
RETURNING user_id, team_id, created_at, updated_at;

-- name: GetTeamMembers :many
SELECT user_id, team_id, created_at, updated_at
FROM users_teams
WHERE team_id = ?1;

-- name: RemoveUserFromTeam :exec
DELETE FROM users_teams
WHERE user_id = ?1 AND team_id = ?2;

-- name: RemoveAllUsersFromTeam :exec
DELETE FROM users_teams
WHERE team_id = ?1;

-- name: IsUserInTeam :one
SELECT EXISTS (
  SELECT 1
  FROM users_teams
  WHERE user_id = ?1 AND team_id = ?2
);