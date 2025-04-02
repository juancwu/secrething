-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vaults (
    vault_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_by_user_id TEXT NOT NULL,  -- Who created it (always a user)
    owner_type TEXT NOT NULL,          -- "user" or "team"
    owner_id TEXT NOT NULL,            -- user_id or team_id
    created_at TEXT NOT NULL CHECK(created_at != ''),
    updated_at TEXT NOT NULL CHECK(updated_at != ''),
    FOREIGN KEY (created_by_user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vaults;
-- +goose StatementEnd
