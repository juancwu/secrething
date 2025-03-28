-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bento_tokens(
    bento_token_id TEXT NOT NULL PRIMARY KEY,
    bento_id TEXT NOT NULL CHECK (bento_id != ''),
    token_salt BLOB NOT NULL UNIQUE,
    created_by TEXT NOT NULL CHECK (created_by != ''),
    created_at TEXT NOT NULL CHECK (created_at != ''),
    last_used_at TEXT,
    expires_at TEXT,

    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(bento_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bento_tokens;
-- +goose StatementEnd
