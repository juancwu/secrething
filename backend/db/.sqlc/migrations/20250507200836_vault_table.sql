-- +goose Up
-- +goose StatementBegin
CREATE TABLE vaults (
    vault_id TEXT NOT NULL PRIMARY KEY,

    vault_name TEXT NOT NULL,
    owner_id TEXT NOT NULL,

    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    schedule_delete TEXT,

    UNIQUE(vault_name, owner_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vaults;
-- +goose StatementEnd
