-- name: NewAuthToken :one
INSERT INTO auth_tokens
(user_id, created_at, expires_at, token_type)
VALUES
(?, ?, ?, ?)
RETURNING *; 

-- name: ExistsAuthTokenById :one
SELECT EXISTS(SELECT 1 FROM auth_tokens WHERE id = ?);

-- name: GetAuthTokenById :one
SELECT * FROM auth_tokens
WHERE id = ?;

-- name: GetUserAuthTokens :many
SELECT * FROM auth_tokens
WHERE user_id = ?;

-- name: DeletAuthTokenById :exec
DELETE FROM auth_tokens WHERE id = ?;

-- name: DeleteUserAuthTokens :exec
DELETE FROM auth_tokens WHERE user_id = ?;

-- name: DeleteAllTokensByTypeAndUserID :exec
DELETE FROM auth_tokens WHERE user_id = ? AND token_type = ?;
