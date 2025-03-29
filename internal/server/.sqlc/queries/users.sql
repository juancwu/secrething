-- name: ExistsUser :one
SELECT 1 FROM users WHERE email = ?1;

-- name: NewUser :one
INSERT INTO users
(
    email,
    password_hash,
    name,
    created_at,
    updated_at
)
VALUES (?1, ?2, ?3, ?4, ?5)
RETURNING *;
