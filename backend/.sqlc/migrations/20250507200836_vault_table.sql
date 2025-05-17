-- +goose Up
-- +goose StatementBegin
CREATE TABLE vaults (
    vault_id TEXT NOT NULL PRIMARY KEY,

    vault_name TEXT NOT NULL,
    vault_owner_id TEXT NOT NULL,

    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    UNIQUE(vault_name, vault_owner_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vaults;
-- +goose StatementEnd
