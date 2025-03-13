-- +goose Up
-- +goose StatementBegin
-- Security notification log
CREATE TABLE security_notifications (
    notification_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    security_event_id TEXT NULL,
    notification_type TEXT NOT NULL, -- 'new_device', 'suspicious_login', 'failed_attempts', etc.
    notification_channel TEXT NOT NULL DEFAULT 'email', -- For future expansion to SMS, push, etc.
    recipient TEXT NOT NULL, -- Email or other contact info
    content TEXT NOT NULL,
    is_sent BOOL NOT NULL DEFAULT false,
    sent_at TEXT NULL,
    error_message TEXT NULL, -- If sending failed

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_security_event_id FOREIGN KEY (security_event_id) REFERENCES security_events(event_id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS security_notifications;
-- +goose StatementEnd
