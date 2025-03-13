-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS totp_recovery_codes (
    totp_recovery_code_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL CHECK (user_id != ''),
    code TEXT NOT NULL CHECK (code != ''),
    used BOOL NOT NULL DEFAULT false,
    created_at TEXT NOT NULL CHECK (created_at != ''),
    used_at TEXT,

    CONSTRAINT unique_user_recovery_code UNIQUE (user_id, code),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS totp_recovery_codes;
-- +goose StatementEnd
