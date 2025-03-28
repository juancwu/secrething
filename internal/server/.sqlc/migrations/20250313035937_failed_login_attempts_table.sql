-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS failed_login_attempts (
    attempt_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NULL, -- Can be NULL if attempt was with non-existent email
    email TEXT NOT NULL, -- Store the email used in the attempt
    ip_address TEXT NOT NULL,
    user_agent TEXT NULL,
    device_fingerprint TEXT NULL,
    geolocation_country TEXT NULL,
    geolocation_city TEXT NULL,
    attempt_time TEXT NOT NULL,
    failure_reason TEXT NOT NULL, -- 'invalid_password', 'invalid_totp', 'account_locked', 'user_not_found'

    -- Optional link to known device if fingerprint matches
    device_id TEXT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL,
    CONSTRAINT fk_device_id FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS failed_login_attempts;
-- +goose StatementEnd
