-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS auth_tokens (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (gen_random_uuid()),
    user_id TEXT NOT NULL CHECK (user_id != ''),
    created_at TEXT NOT NULL CHECK (created_at != ''),
    expires_at TEXT NOT NULL CHECK (expires_at != ''),
    token_type TEXT NOT NULL CHECK (token_type != ''),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth_tokens;
-- +goose StatementEnd
