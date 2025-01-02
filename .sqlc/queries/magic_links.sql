-- name: CreateMagicLink :exec
INSERT INTO magic_links
(token, user_id, created_at, expires_at)
VALUES
(?, ?, ?, ?);

-- name: GetMagicLink :one
SELECT token, user_id, created_at, expires_at
FROM magic_links
WHERE token = ? AND user_id = ?;

-- name: RemoveMagicLink :exec
DELETE FROM magic_links
WHERE token = ? AND user_id = ?;
