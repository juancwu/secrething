-- name: NewJWT :one
INSERT INTO jwts
(user_id, created_at, expires_at, token_type)
VALUES
(?, ?, ?, ?)
RETURNING *; 

-- name: ExistsJwtById :one
SELECT EXISTS(SELECT 1 FROM jwts WHERE id = ?);

-- name: GetJwtById :one
SELECT * FROM jwts
WHERE id = ?;

-- name: GetUserJwts :many
SELECT * FROM jwts
WHERE user_id = ?;

-- name: DeletJwtById :exec
DELETE FROM jwts WHERE id = ?;

-- name: DeleteUserJwts :exec
DELETE FROM jwts WHERE user_id = ?;
