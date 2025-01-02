-- name: CreateMagicLink :exec
INSERT INTO magic_links
(token, user_id, created_at, expires_at)
VALUES
(?, ?, ?, ?);

-- name: GetMagicLink :one
SELECT ml.token, ml.user_id, ml.created_at, ml.expires_at
        FROM magic_links AS ml
    LEFT JOIN users AS u ON u.email = ?
    WHERE ml.token = ?;

-- name: RemoveMagicLink :exec
DELETE FROM magic_links
WHERE token = ? AND user_id = ?;
