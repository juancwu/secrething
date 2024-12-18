CREATE TABLE email_tokens (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid4()),
    user_id TEXT NOT NULL UNIQUE,
    token_salt BLOB NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),
    expires_at TEXT NOT NULL
);
