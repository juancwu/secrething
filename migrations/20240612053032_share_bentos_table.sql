-- +goose Up
-- +goose StatementBegin
CREATE TABLE shared_bentos (
    id SERIAL PRIMARY KEY,
    owner_id UUID NOT NULL,
    recipient_id UUID NOT NULL,
    shared_bento_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_shared_bento UNIQUE (owner_id, shared_bento_id, recipient_id),
    CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_recipient_id FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_shared_bento_id FOREIGN KEY (shared_bento_id) REFERENCES bentos(id) ON DELETE CASCADE
);

-- this is the auto update the updated_at columnd
CREATE TRIGGER update_shared_bentos_updated_at
BEFORE UPDATE ON shared_bentos
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column(); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_shared_bentos_updated_at ON shared_bentos;
DROP TABLE shared_bentos;
-- +goose StatementEnd
