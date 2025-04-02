-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_tokens (
    user_token_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL CHECK(user_id != ''),
    token_type TEXT NOT NULL,
    expires_at TEXT NOT NULL CHECK(expires_at != ''),
    created_at TEXT NOT NULL CHECK(created_at != ''),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_tokens;
-- +goose StatementEnd
