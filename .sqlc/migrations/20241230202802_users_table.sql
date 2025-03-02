-- +goose Up
-- +goose StatementBegin

-- enable foreign key support
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE CHECK (email != ''),
    password TEXT NOT NULL CHECK (password != ''),
    nickname TEXT NOT NULL CHECK (nickname != ''),
    email_verified BOOL NOT NULL DEFAULT false,
    totp_secret TEXT,
    totp_locked BOOL NOT NULL DEFAULT false,
    created_at TEXT NOT NULL CHECK(created_at != ''),
    updated_at TEXT NOT NULL CHECK(updated_at != '')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- disable foreign key support
PRAGMA foreign_keys = OFF;
-- +goose StatementEnd
