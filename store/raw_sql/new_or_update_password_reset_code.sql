INSERT INTO password_reset_codes (user_id, code, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE SET
    code = EXCLUDED.code,
    expires_at = EXCLUDED.expires_at
RETURNING id, created_at, updated_at;
