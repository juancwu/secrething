-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS environments (
    environment_id TEXT NOT NULL PRIMARY KEY,
    vault_id TEXT NOT NULL CHECK (vault_id != ''),
    name TEXT NOT NULL CHECK (name != ''),
    description TEXT,
    created_at TEXT NOT NULL CHECK (created_at != ''),
    updated_at TEXT NOT NULL CHECK (updated_at != ''),
    CONSTRAINT unique_environment_name_vault UNIQUE (vault_id, name),
    CONSTRAINT fk_vault_id FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS environments;
-- +goose StatementEnd