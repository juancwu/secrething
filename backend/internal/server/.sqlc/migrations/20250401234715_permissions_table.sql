-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS permissions (
    permission_id TEXT NOT NULL PRIMARY KEY,
    vault_id TEXT NOT NULL CHECK(vault_id != ''),
    grantee_type TEXT NOT NULL CHECK(grantee_type != ''),        -- "user" or "team"
    grantee_id TEXT NOT NULL CHECK(grantee_id != ''),          -- user_id or team_id
    permission_bits BIGINT NOT NULL,   -- Bitmask for granular permissions
    granted_by TEXT NOT NULL CHECK(granted_by != ''),          -- user_id who granted access
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES users(user_id),
    UNIQUE (vault_id, grantee_type, grantee_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS permissions;
-- +goose StatementEnd
