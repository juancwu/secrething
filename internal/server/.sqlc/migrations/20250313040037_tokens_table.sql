-- +goose Up
-- +goose StatementBegin
CREATE TABLE tokens (
    token_id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    device_id TEXT NULL,
    token_hash TEXT NOT NULL UNIQUE, -- Store hash of the token, not the token itself
    token_type TEXT NOT NULL,
    access_level TEXT NOT NULL,
    is_active BOOL NOT NULL DEFAULT false,
    issued_at TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    last_used_at TEXT NOT NULL,
    revoked_at TEXT NULL,
    revocation_reason TEXT NULL,

    -- For refresh tokens, link to the access token they're associated with
    parent_token_id TEXT NULL,

    -- IP and context when token was created
    issued_ip TEXT NULL,
    issued_location_country TEXT NULL,
    issued_location_city TEXT NULL,

    -- Additional security context
    session_id TEXT NULL, -- To group related tokens into a session
    scope TEXT NULL, -- For granular permission control (comma-separated list)

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_device_id FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE,
    CONSTRAINT fk_parent_token_id FOREIGN KEY (parent_token_id) REFERENCES tokens(token_id) ON DELETE CASCADE
);

-- Create index for token lookup
CREATE INDEX idx_tokens_token_hash ON tokens(token_hash);
CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_device_id ON tokens(device_id);
CREATE INDEX idx_active_tokens ON tokens(is_active, expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tokens_token_hash;
DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_device_id;
DROP INDEX IF EXISTS idx_active_tokens;

DROP TABLE IF EXISTS tokens;
-- +goose StatementEnd
