-- name: CreateEmailToken :one
INSERT INTO email_tokens
(user_id, token_salt, created_at, expires_at)
VALUES (?, ?, ?, ?)
RETURNING id;

-- name: DeleteAllEmailTokensByUserId :many
DELETE FROM email_tokens
WHERE user_id = ?
RETURNING id;

-- name: GetEmailTokenById :one
SELECT
    id,
    user_id,
    token_salt,
    created_at,
    expires_at
FROM email_tokens
WHERE id = ?;
