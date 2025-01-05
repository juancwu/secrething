-- +goose Up
-- +goose StatementBegin

-- enable foreign key support
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL PRIMARY KEY DEFAULT(gen_random_uuid()),
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    nickname TEXT NOT NULL,
    email_verified BOOL NOT NULL DEFAULT false,
    totp_secret TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- disable foreign key support
PRAGMA foreign_keys = OFF;
-- +goose StatementEnd
