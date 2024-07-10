-- +goose Up
-- +goose StatementBegin
CREATE TABLE bento_entries (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    bento_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_bento_entry UNIQUE (name, bento_id),
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(id) ON DELETE CASCADE
);
CREATE TRIGGER update_bento_entries_updated_at
BEFORE UPDATE ON bento_entries
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column(); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_bento_entries_updated_at ON bento_entries;
DRO; TABLE bento_entries;
-- +goose StatementEnd
