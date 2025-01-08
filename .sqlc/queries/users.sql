-- name: CreateUser :one
INSERT INTO users
(email, password, nickname, created_at, updated_at)
VALUES
(?, ?, ?, ?, ?)
RETURNING id;

-- name: ExistsUserWithEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = ?);

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = ?;

-- name: SetUserEmailVerifiedStatus :exec
UPDATE users SET email_verified = ?, updated_at = ? WHERE id = ?;
