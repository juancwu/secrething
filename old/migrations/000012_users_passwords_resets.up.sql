CREATE TABLE users_passwords_resets (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    reset_id CHAR(12) NOT NULL,

    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- no updated_at because this entry should be deleted right after it has been used

    -- foreign keys
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id),

    -- unique
    CONSTRAINT unique_reset_id UNIQUE (reset_id)
);
