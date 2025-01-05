-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS partial_tokens (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (gen_random_uuid()),
    user_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    expires_at TEXT NO NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS partial_tokens;
-- +goose StatementEnd
