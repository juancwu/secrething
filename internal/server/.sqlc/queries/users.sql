-- name: UserExistsWithEmail :one
SELECT email FROM users WHERE email = ?1;
