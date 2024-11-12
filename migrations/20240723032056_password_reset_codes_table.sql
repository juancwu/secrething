-- +goose Up
-- +goose StatementBegin
CREATE TABLE password_reset_codes (
    id SERIAL NOT NULL PRIMARY KEY,
    code CHAR(6) NOT NULL,
    user_id UUID NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TRIGGER update_password_reset_codes_updated_at
BEFORE UPDATE ON password_reset_codes
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column(); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_password_reset_codes_updated_at ON password_reset_codes;
DROP TABLE password_reset_codes;
-- +goose StatementEnd
