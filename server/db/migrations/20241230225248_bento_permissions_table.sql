-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bento_permissions (
    user_id TEXT NOT NULL CHECK (user_id != ''),
    bento_id TEXT NOT NULL CHECK (bento_id != ''),
    bytes BLOB NOT NULL,
    created_at TEXT NOT NULL CHECK (created_at != ''),
    updated_at TEXT NOT NULL CHECK (updated_at != ''),
    CONSTRAINT pk_bento_permissions PRIMARY KEY (user_id, bento_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bento_permissions;
-- +goose StatementEnd
