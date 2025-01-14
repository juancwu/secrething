-- name: NewBentoPermission :exec
INSERT INTO bento_permissions
(user_id, bento_id, bytes, created_at, updated_at)
VALUES
(?, ?, ?, ?, ?);
