-- +goose Up
-- +goose StatementBegin
-- User verification table to track verification flow
CREATE TABLE user_verifications (
    verification_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    verification_type TEXT NOT NULL, -- 'email', 'totp_setup', 'phone', etc.
    verification_token_id TEXT NULL,
    verification_code TEXT NULL, -- For email verification codes
    is_verified BOOL NOT NULL DEFAULT false,
    attempts INTEGER NOT NULL DEFAULT 0,
    expires_at TEXT NOT NULL,
    verified_at TEXT NULL,
    created_at TEXT NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_verification_token_id FOREIGN KEY (verification_token_id) REFERENCES tokens(token_id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_verifications;
-- +goose StatementEnd
