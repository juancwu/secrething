-- +goose Up
-- +goose StatementBegin
CREATE TABLE bento_entries (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    content TEXT NOT NULL,
    bento_id UUID NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(id) ON DELETE CASCADE,

    CONSTRAINT unique_name_bento_entry UNIQUE (name, bento_id)
);

-- trigger to update the updated_at column
CREATE TRIGGER update_bento_entries_updated_at
BEFORE UPDATE ON bento_entries
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bento_entries;
DROP TRIGGER IF EXISTS update_bento_entries_updated_at ON bento_entries;
-- +goose StatementEnd
