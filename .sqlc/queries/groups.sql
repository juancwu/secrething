-- name: NewGroup :one
INSERT INTO groups (name, owner_id, created_at, updated_at)
VALUES (?, ?, ?, ?)
RETURNING id;

-- name: AddUserToGroup :exec
INSERT INTO users_groups (user_id, group_id, created_at) VALUES (?, ?, ?);

-- name: RemoveUserFromGroup :exec
DELETE FROM users_groups WHERE user_id = ? AND group_id = ?;

-- name: ExistsGroupOwnedByUser :one
SELECT EXISTS(SELECT 1 FROM groups WHERE name = ? AND owner_id = ?);

-- name: ExistsGroupWithIdOwnedByUser :one
SELECT EXISTS(SELECT 1 FROM groups WHERE id = ? AND owner_id = ?);

-- name: RemoveGroupByID :exec
DELETE FROM groups WHERE id = ?;
