SELECT
    id,
    user_id,
    bento_id,
    permissions,
    created_at,
    updated_at
FROM bento_permissions WHERE user_id = $1 AND bento_id = $2;
