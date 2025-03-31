-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS team_invitations (
    team_invitation_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL CHECK (user_id != ''),
    team_id TEXT NOT NULL CHECK (team_id != ''),
    resend_email_id TEXT NOT NULL,
    created_at TEXT NOT NULL CHECK (created_at != ''),
    expires_at TEXT NOT NULL CHECK (expires_at != ''),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_team_id FOREIGN KEY (team_id) REFERENCES teams(team_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team_invitations;
-- +goose StatementEnd
