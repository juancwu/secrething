-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vaults (
    vault_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL CHECK (user_id != ''),
    name TEXT NOT NULL CHECK (name != ''),
    created_at TEXT NOT NULL CHECK (created_at != ''),
    updated_at TEXT NOT NULL CHECK (updated_at != ''),
    CONSTRAINT unique_vault_name_user UNIQUE (user_id, name),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vaults;
-- +goose StatementEnd
