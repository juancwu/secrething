-- name: NewPartialToken :one
INSERT INTO partial_tokens
(user_id, created_at, updated_at, expires_at)
VALUES
(?, ?, ?, ?)
RETURNING id;

-- name: GetPartialTokenById :one
SELECT
    id,
    user_id,
    created_at,
    updated_at,
    expires_at
FROM partial_tokens
WHERE id = ?;

-- name: GetUserPartialTokens :many
SELECT
    id,
    user_id,
    created_at,
    updated_at,
    expires_at
FROM partial_tokens
WHERE user_id = ?;

-- name: DeletePartialTokenById :exec
DELETE FROM partial_tokens WHERE id = ?;

-- name: DeleteUserPartialTokens :exec
DELETE FROM partial_tokens WHERE user_id = ?;
