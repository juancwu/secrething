-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_teams (
    user_id TEXT NOT NULL CHECK(user_id != ''),
    team_id TEXT NOT NULL CHECK(team_id != ''),
    created_at TEXT NOT NULL CHECK(created_at != ''),
    updated_at TEXT NOT NULL CHECK(updated_at != ''),
    PRIMARY KEY (user_id, team_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(team_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_teams;
-- +goose StatementEnd
