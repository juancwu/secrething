-- +goose Up
-- +goose StatementBegin
-- Token usage/audit log for security monitoring
CREATE TABLE IF NOT EXISTS token_events (
    event_id TEXT NOT NULL PRIMARY KEY,
    token_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    device_id TEXT NULL,
    event_type TEXT NOT NULL, -- 'used', 'refreshed', 'expired', 'revoked'
    ip_address TEXT NULL,
    user_agent TEXT NULL,
    event_details_json TEXT NULL, -- JSON as TEXT in SQLite
    created_at TEXT NOT NULL,

    CONSTRAINT fk_token_id FOREIGN KEY (token_id) REFERENCES tokens(token_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_device_id FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS token_events;
-- +goose StatementEnd
