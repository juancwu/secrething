-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_settings (
    user_setting_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,

    -- Settings
    alert_on_suspicious_activity BOOL NOT NULL DEFAULT true,

    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_settings;
-- +goose StatementEnd
