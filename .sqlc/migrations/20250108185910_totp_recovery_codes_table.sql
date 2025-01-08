-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS totp_recovery_codes (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (gen_random_uuid()),
    user_id TEXT NOT NULL,
    code TEXT NOT NULL,
    used BOOL NOT NULL DEFAULT false,
    created_at TEXT NOT NULL,

    CONSTRAINT unique_user_recovery_code UNIQUE (user_id, code),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS totp_recovery_codes;
-- +goose StatementEnd
