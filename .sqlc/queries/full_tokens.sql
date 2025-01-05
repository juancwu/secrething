-- name: NewFullToken :one
INSERT INTO full_tokens
(user_id, created_at, updated_at, expires_at)
VALUES
(?, ?, ?, ?)
RETURNING id;

-- name: GetFullTokenById :one
SELECT
    id,
    user_id,
    created_at,
    updated_at,
    expires_at
FROM full_tokens
WHERE id = ?;

-- name: GetUserFullTokens :many
SELECT
    id,
    user_id,
    created_at,
    updated_at,
    expires_at
FROM full_tokens
WHERE user_id = ?;

-- name: DeleteFullTokenById :exec
DELETE FROM full_tokens WHERE id = ?;

-- name: DeleteUserFullTokens :exec
DELETE FROM full_tokens WHERE user_id = ?;
