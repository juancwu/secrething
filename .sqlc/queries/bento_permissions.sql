-- name: NewBentoPermission :exec
INSERT INTO bento_permissions
(user_id, bento_id, bytes, created_at, updated_at)
VALUES
(?, ?, ?, ?, ?);

-- name: GetUserIDsWithBentoAccess :many
SELECT user_id FROM bento_permissions WHERE bento_id = ?;
