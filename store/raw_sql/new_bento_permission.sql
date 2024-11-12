INSERT INTO bento_permissions (user_id, bento_id, permissions)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at;
