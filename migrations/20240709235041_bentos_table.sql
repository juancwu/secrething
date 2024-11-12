-- +goose Up
-- +goose StatementBegin
CREATE TABLE bentos (
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    owner_id UUID NOT NULL,
    pub_key TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_bento_owner UNIQUE (name, owner_id),
    CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TRIGGER update_bentos_updated_at
BEFORE UPDATE ON bentos
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column(); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_bentos_updated_at ON bentos;
DROP TABLE bentos;
-- +goose StatementEnd
