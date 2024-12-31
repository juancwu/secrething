-- name: CreateEmailToken :one
INSERT INTO email_tokens
(user_id, token_salt, created_at, expires_at)
VALUES
()
