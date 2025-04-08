-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
    team_id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_by_user_id TEXT NOT NULL CHECK(created_by_user_id != ''),
    created_at TEXT NOT NULL CHECK(created_at != ''),
    updated_at TEXT NOT NULL CHECK(updated_at != ''),
    FOREIGN KEY (created_by_user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
