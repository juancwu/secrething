CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY CHECK (length(id) = 36), -- UUID format
    email TEXT NOT NULL UNIQUE COLLATE NOCASE,
    password_hash TEXT NOT NULL CHECK (length(password_hash) > 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    last_login_at TEXT,

    CHECK (email LIKE '%@%.%') -- Basic email format validation
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id TEXT PRIMARY KEY CHECK (length(id) = 36), -- UUID format
    user_id TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- create indices
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- create trigger to update the updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_users_timestamp
    AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = datetime('now')
    WHERE id = NEW.id;
END;

-- create trigger to clean up refresh tokens when user is deleted
CREATE TRIGGER IF NOT EXISTS cleanup_refresh_tokens_after_user_delete
    AFTER DELETE ON users
BEGIN
    DELETE FROM refresh_tokens WHERE user_id = OLD.id;
END;
