-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_verifications (
    id SERIAL NOT NULL PRIMARY KEY,
    code CHAR(20) NOT NULL,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_user_code UNIQUE (code, user_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_verifications;
-- +goose StatementEnd
