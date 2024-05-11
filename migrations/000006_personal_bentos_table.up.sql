-- personal bentos
CREATE TABLE personal_bentos (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    owner_id UUID,
    content BYTEA NOT NULL, -- encrypted
    pub_key BYTEA NOT NULL, -- encrypted

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_bento_owner UNIQUE (name, owner_id),
    CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);
