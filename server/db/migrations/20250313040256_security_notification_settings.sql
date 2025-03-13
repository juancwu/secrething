-- +goose Up
-- +goose StatementBegin
CREATE TABLE security_notification_settings (
    user_id TEXT NOT NULL PRIMARY KEY,
    notify_on_new_device BOOL NOT NULL DEFAULT true,
    notify_on_suspicious_login BOOL NOT NULL DEFAULT true,
    notify_on_failed_attempts INTEGER NOT NULL DEFAULT 1, -- Only after threshold
    failed_attempts_threshold INTEGER NOT NULL DEFAULT 3, -- Send notification after X failed attempts
    notify_on_password_change BOOL NOT NULL DEFAULT true,
    notify_on_email_change BOOL NOT NULL DEFAULT true,
    notify_on_totp_change BOOL NOT NULL DEFAULT true,
    notification_email TEXT NULL, -- Separate email for security notifications

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS security_notification_settings;
-- +goose StatementEnd
