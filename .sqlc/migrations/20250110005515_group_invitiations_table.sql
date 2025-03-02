-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS group_invitations (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL CHECK (user_id != ''),
    group_id TEXT NOT NULL CHECK (group_id != ''),
    created_at TEXT NOT NULL CHECK (created_at != ''),
    expires_at TEXT NOT NULL CHECK (expires_at != ''),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS group_invitations;
-- +goose StatementEnd
