-- name: NewJWT :one
INSERT INTO jwts
(user_id, created_at, expires_at, token_type)
VALUES
(?, ?, ?, ?)
RETURNING *; 

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
