-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS devices (
    device_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    device_name TEXT NULL,
    device_type TEXT NULL, -- mobile, desktop, tablet, etc.
    os_type TEXT NULL,
    os_version TEXT NULL,
    app_version TEXT NULL,
    browser_type TEXT NULL,
    browser_version TEXT NULL,
    ip_address TEXT NULL,
    user_agent TEXT NULL,
    device_fingerprint TEXT NULL,
    is_trusted BOOL DEFAULT false,
    first_seen_at TEXT NULL,
    last_active_at TEXT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,

    -- Optional geographical data for advanced security monitoring
    last_latitude REAL NULL,
    last_longitude REAL NULL,
    last_location_country TEXT NULL,
    last_location_city TEXT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS devices;
-- +goose StatementEnd
