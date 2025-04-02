-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE CHECK (email != ''),
    password_hash TEXT NOT NULL CHECK (password_hash != ''),
    name TEXT,
    email_verified BOOL NOT NULL DEFAULT false,
    totp_secret TEXT,
    totp_enabled BOOL NOT NULL DEFAULT false,
    account_status TEXT NOT NULL DEFAULT 'pending', -- pending, active, suspended, locked
    failed_login_attempts INTEGER DEFAULT 0, -- Kept for quick access but detailed tracking in separate table, failed_login_attempts
    last_failed_login_at TEXT,
    account_locked_until TEXT,
    created_at TEXT NOT NULL CHECK(created_at != ''),
    updated_at TEXT NOT NULL CHECK(updated_at != '')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
