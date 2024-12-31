-- name: CreateUser :one
INSERT INTO users
(email, password, nickname, token_salt, created_at, updated_at)
VALUES
(?, ?, ?, ?, ?, ?)
RETURNING id, email_verified;
