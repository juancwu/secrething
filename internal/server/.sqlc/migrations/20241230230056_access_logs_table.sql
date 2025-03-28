-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access_logs (
    access_log_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT,
    bento_id TEXT,
    group_id TEXT,
    bento_token_id TEXT,
    action TEXT NOT NULL,
    details JSONB,
    accessed_at TEXT NOT NULL CHECK (accessed_at != '')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_logs;
-- +goose StatementEnd
