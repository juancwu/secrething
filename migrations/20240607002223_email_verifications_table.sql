-- +goose Up
-- +goose StatementBegin

CREATE TABLE email_verifications (
    id SERIAL PRIMARY KEY NOT NULL,
    code CHAR(16) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- foreign keys
    CONSTRAINT fk_email_verification_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER update_email_verifications_updated_at
BEFORE UPDATE ON email_verifications
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_email_verifications_updated_at ON email_verifications;
DROP TABLE email_verifications;
-- +goose StatementEnd
