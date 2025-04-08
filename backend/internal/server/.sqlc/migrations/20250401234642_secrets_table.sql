-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS secrets (
    secret_id TEXT NOT NULL PRIMARY KEY,
    vault_id TEXT NOT NULL CHECK (vault_id != ''),
    name TEXT NOT NULL CHECK (name != ''),
    value BLOB NOT NULL,
    created_by_user_id TEXT,
    created_at TEXT NOT NULL CHECK(created_at != ''),
    updated_at TEXT NOT NULL CHECK(updated_at != ''),
    FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE,
    FOREIGN KEY (created_by_user_id) REFERENCES users(user_id) ON DELETE SET NULL,
    UNIQUE (vault_id, name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS secrets;
-- +goose StatementEnd
