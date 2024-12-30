-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access_logs (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid4()),
    user_id TEXT,
    bento_id TEXT,
    group_id TEXT,
    bento_token_id TEXT,
    action TEXT NOT NULL,
    details JSONB,
    accessed_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_logs;
-- +goose StatementEnd
