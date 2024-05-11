CREATE TABLE personal_bento_access (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    bento_id UUID NOT NULL,
    can_write BOOLEAN NOT NULL DEFAULT false,
    can_read BOOLEAN NOT NULL DEFAULT false,
    can_delete BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES personal_bentos(id) ON DELETE CASCADE,

    CONSTRAINT unique_user_bento_access UNIQUE (user_id, bento_id)
);
