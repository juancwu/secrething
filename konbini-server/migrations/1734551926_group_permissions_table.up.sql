CREATE TABLE group_permissions(
    group_id TEXT NOT NULL,
    bento_id TEXT NOT NULL,
    level INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),

    CONSTRAINT pk_group_permissions PRIMARY KEY (group_id, bento_id)
);
