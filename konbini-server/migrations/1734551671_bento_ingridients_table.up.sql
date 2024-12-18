CREATE TABLE bento_ingridients (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid4()),
    bento_id TEXT NOT NULL,
    name TEXT NOT NULL,
    value BLOB NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),

    CONSTRAINT unique_bento_ingridient_name UNIQUE (bento_id, name)
);
