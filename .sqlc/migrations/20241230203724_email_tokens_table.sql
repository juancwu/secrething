-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS email_tokens (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (gen_random_uuid()),
    user_id TEXT NOT NULL UNIQUE,
    token_salt BLOB NOT NULL UNIQUE,
    created_at TEXT NOT NULL,
    expires_at TEXT NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS email_tokens;
-- +goose StatementEnd
