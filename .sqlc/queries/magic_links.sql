-- name: CreateMagicLink :one
INSERT INTO magic_links
(user_id, state, created_at, expires_at)
VALUES
(?, ?, ?, ?)
RETURNING id;

-- name: GetMagicLink :one
SELECT id, user_id, state, created_at, expires_at
FROM magic_links
WHERE id = ? AND user_id = ? AND state = ?;

-- name: RemoveMagicLink :exec
DELETE FROM magic_links
WHERE id = ? AND user_id = ? AND state = ?;
