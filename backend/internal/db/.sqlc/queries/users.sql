-- name: GetUserByEmail :one
SELECT
    user_id,
    email,
    first_name,
    last_name,
    created_at,
    updated_at
FROM users WHERE email = ?1;

-- name: GetUserByID :one
SELECT
    user_id,
    email,
    first_name,
    last_name,
    created_at,
    updated_at
FROM users WHERE user_id = ?1;
