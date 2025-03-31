-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vault_tokens(
    vault_token_id TEXT NOT NULL PRIMARY KEY,
    vault_id TEXT NOT NULL CHECK (vault_id != ''),
    token_salt BLOB NOT NULL UNIQUE,
    created_by TEXT NOT NULL CHECK (created_by != ''),
    created_at TEXT NOT NULL CHECK (created_at != ''),
    last_used_at TEXT,
    expires_at TEXT,

    CONSTRAINT fk_vault_id FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vault_tokens;
-- +goose StatementEnd
