-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_groups (
    user_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    CONSTRAINT pk_users_groups PRIMARY KEY (user_id, group_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_groups;
-- +goose StatementEnd
