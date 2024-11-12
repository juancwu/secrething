-- drop triggers first
DROP TRIGGER IF EXISTS update_users_timestamp;
DROP TRIGGER IF EXISTS cleanup_refresh_tokens_after_user_delete;

-- drop indices
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;

-- drop tables
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
