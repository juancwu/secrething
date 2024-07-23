SELECT
    id,
    code,
    user_id,
    expires_at,
    created_at,
    updated_at
FROM password_reset_codes WHERE user_id = $1;
