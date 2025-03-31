-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
    team_id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL CHECK (name != ''),
    owner_id TEXT NOT NULL CHECK (owner_id != ''),
    created_at TEXT NOT NULL CHECK (created_at != ''),
    updated_at TEXT NOT NULL CHECK (updated_at != ''),

    CONSTRAINT unique_team_name_owner UNIQUE (name, owner_id),

    CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
