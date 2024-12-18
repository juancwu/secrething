CREATE TABLE users (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid4()),
    email TEXT NOT NULL UNIQUE,
    password BLOB NOT NULL,
    nickname TEXT NOT NULL,
    email_verified BOOL NOT NULL DEFAULT false,
    token_salt BLOB NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'utc'))
);
