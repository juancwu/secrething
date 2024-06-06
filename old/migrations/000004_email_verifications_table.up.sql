CREATE TYPE email_status AS ENUM (
    'PENDING',
    'SENT',
    'OPENED',
    'VERIFIED',
    'FAILED'
);

CREATE TABLE email_verifications (
    id SERIAL PRIMARY KEY NOT NULL,
    verification_id CHAR(16) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    resend_email_id TEXT,
    status EMAIL_STATUS NOT NULL DEFAULT 'PENDING',

    email_sent_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    verified_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- foreign keys
    CONSTRAINT fk_email_verification_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
