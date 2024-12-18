CREATE TABLE groups (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid4()),
    name TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),

    CONSTRAINT unique_group_name_owner UNIQUE (name, owner_id)
);
