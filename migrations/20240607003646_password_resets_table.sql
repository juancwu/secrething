-- +goose Up
-- +goose StatementBegin
CREATE TABLE password_resets (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    reset_code CHAR(6) NOT NULL,

    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- no updated_at because this entry should be deleted right after it has been used

    -- foreign keys
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id),

    -- unique
    CONSTRAINT unique_reset_code UNIQUE (reset_code)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE password_resets;
-- +goose StatementEnd
