CREATE TABLE access_logs (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid4()),
    user_id TEXT,
    bento_id TEXT,
    group_id TEXT,
    bento_token_id TEXT,
    action TEXT NOT NULL,
    details JSONB,
    accessed_at TEXT NOT NULL DEFAULT (datetime('now', 'utc'))
);
