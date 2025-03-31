-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS team_permissions(
    team_id TEXT NOT NULL CHECK (team_id != ''),
    vault_id TEXT NOT NULL CHECK (vault_id != ''),
    permissions INTEGER NOT NULL,
    created_at TEXT NOT NULL CHECK (created_at != ''),
    updated_at TEXT NOT NULL CHECK (updated_at != ''),
    CONSTRAINT pk_team_permissions PRIMARY KEY (team_id, vault_id),
    CONSTRAINT fk_team_id FOREIGN KEY (team_id) REFERENCES teams(team_id) ON DELETE CASCADE,
    CONSTRAINT fk_vault_id FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team_permissions;
-- +goose StatementEnd
