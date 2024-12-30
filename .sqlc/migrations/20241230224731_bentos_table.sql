-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bentos (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid4()),
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    CONSTRAINT unique_bento_name_user UNIQUE (user_id, name),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bentos;
-- +goose StatementEnd
