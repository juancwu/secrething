-- name: CreateSession :one
INSERT INTO sessions
(token_salt, user_id, ip, last_activity)
VALUES
(?, ?, ?, ?)
RETURNING token_id;
