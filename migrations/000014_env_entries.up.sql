CREATE TABLE personal_bento_entries (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    content TEXT NOT NULL,
    personal_bento_id UUID NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_personal_bento_id FOREIGN KEY (personal_bento_id) REFERENCES personal_bentos(id) ON DELETE CASCADE,

    CONSTRAINT unique_name_personal_bento_entry UNIQUE (name, personal_bento_id)
);

-- trigger to update the updated_at column
CREATE TRIGGER update_personal_bento_entries_updated_at
BEFORE UPDATE ON personal_bento_entries
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
