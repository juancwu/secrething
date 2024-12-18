CREATE TABLE bento_tokens(
    id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid4()),
    bento_id TEXT NOT NULL,
    token_salt BLOB NOT NULL UNIQUE,
    created_by TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),
    last_used_at TEXT,
    expires_at TEXT
);
