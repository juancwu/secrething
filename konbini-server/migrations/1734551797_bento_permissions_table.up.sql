CREATE TABLE bento_permissions (
    user_id TEXT NOT NULL,
    bento_id TEXT NOT NULL,
    -- default to no permissions
    level INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),

    CONSTRAINT pk_bento_permissions PRIMARY KEY (user_id, bento_id)
);
