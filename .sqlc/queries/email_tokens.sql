-- name: CreateEmailToken :one
INSERT INTO email_tokens
(user_id, token_salt, created_at, expires_at)
VALUES (?, ?, ?, ?)
RETURNING id;

-- name: DeleteAllEmailTokensByUserId :many
DELETE FROM email_tokens
WHERE user_id = ?
RETURNING id;
