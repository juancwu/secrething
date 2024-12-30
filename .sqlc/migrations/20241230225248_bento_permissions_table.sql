-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bento_permissions (
    user_id TEXT NOT NULL,
    bento_id TEXT NOT NULL,
    -- default to no permissions
    level INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    CONSTRAINT pk_bento_permissions PRIMARY KEY (user_id, bento_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bento_permissions;
-- +goose StatementEnd
