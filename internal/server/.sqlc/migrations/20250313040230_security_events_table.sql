-- +goose Up
-- +goose StatementBegin
-- Security events table for suspicious activities
CREATE TABLE security_events (
    event_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NULL, -- NULL if event is not tied to a specific user (e.g., IP-based attack)
    device_id TEXT NULL,
    token_id TEXT NULL,
    ip_address TEXT NULL,
    device_fingerprint TEXT NULL,
    event_type TEXT NOT NULL, -- 'suspicious_login', 'token_reuse', 'multiple_locations', 'brute_force', etc.
    severity TEXT NOT NULL, -- 'low', 'medium', 'high', 'critical'
    details_json TEXT NOT NULL, -- JSON as TEXT in SQLite
    is_resolved BOOL NOT NULL DEFAULT false,
    alert_sent BOOL NOT NULL DEFAULT false, -- Boolean to track if email alert was sent
    alert_sent_at TEXT NULL,
    resolution_notes TEXT NULL,
    created_at TEXT NOT NULL,
    resolved_at TEXT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL,
    CONSTRAINT fk_device_id FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE SET NULL,
    CONSTRAINT fk_token_id FOREIGN KEY (token_id) REFERENCES tokens(token_id) ON DELETE SET NULL
);

CREATE INDEX idx_security_events_user ON security_events(user_id, created_at);
CREATE INDEX idx_security_events_ip ON security_events(ip_address, created_at);
CREATE INDEX idx_security_events_type ON security_events(event_type, severity, is_resolved);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_security_events_user;
DROP INDEX IF EXISTS idx_security_events_ip;
DROP INDEX IF EXISTS idx_security_events_type;

DROP TABLE IF EXISTS security_events;
-- +goose StatementEnd
