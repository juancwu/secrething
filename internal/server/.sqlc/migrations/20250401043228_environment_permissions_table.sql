-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS environment_permissions(
    environment_permission_id TEXT NOT NULL PRIMARY KEY,
    vault_id TEXT NOT NULL CHECK (vault_id != ''),
    environment_id TEXT NOT NULL CHECK (environment_id != ''),
    user_id TEXT NOT NULL CHECK (user_id != ''),
    permissions INTEGER NOT NULL,
    created_at TEXT NOT NULL CHECK (created_at != ''),
    updated_at TEXT NOT NULL CHECK (updated_at != ''),
    CONSTRAINT fk_environment_id FOREIGN KEY (environment_id) REFERENCES environments(environment_id) ON DELETE CASCADE,
    CONSTRAINT fk_vault_id FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS environment_permissions;
-- +goose StatementEnd
