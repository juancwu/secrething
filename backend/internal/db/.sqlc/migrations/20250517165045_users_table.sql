-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    user_id TEXT NOT NULL PRIMARY KEY,

    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,

    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX idx_users_email ON users(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users_email;
DROP TABLE users;
-- +goose StatementEnd
