-- +goose Up
-- +goose StatementBegin
CREATE TABLE challenges (
    id SERIAL PRIMARY KEY,
    state CHAR(32) NOT NULL UNIQUE,
    value CHAR(43) NOT NULL,
    user_id UUID NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT fk_challenge_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE challengs;
-- +goose StatementEnd
