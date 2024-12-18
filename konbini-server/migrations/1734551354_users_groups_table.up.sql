CREATE TABLE users_groups (
    user_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),

    CONSTRAINT pk_users_groups PRIMARY KEY (user_id, group_id)
);
