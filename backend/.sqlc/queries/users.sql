-- name: GetUserByEmail :one
SELECT
    user_id,
    email,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at
FROM users WHERE email = ?1;

-- name: GetUserByID :one
SELECT
    user_id,
    email,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at
FROM users WHERE user_id = ?1;

-- name: CreateUser :one
INSERT INTO users (
    user_id,
    email,
    password_hash,
    first_name,
    last_name,
    created_at,
    updated_at
) VALUES (
    ?1, ?2, ?3, ?4, ?5, ?6, ?7
)
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = ?2, updated_at = ?3
WHERE user_id = ?1;