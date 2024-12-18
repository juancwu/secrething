CREATE TABLE sessions (
    token_id TEXT NOT NULL,
    token_salt BLOB NOT NULL UNIQUE,
    user_id TEXT NOT NULL,
    device_name TEXT,
    device_os TEXT,
    device_hostname TEXT,
    ip TEXT,
    location TEXT,
    last_activity TEXT NOT NULL,

    -- define compound pk for sessions since each user should only have one session
    -- per token, if it isn't then there is something fishy going on.
    CONSTRAINT pk_sessions PRIMARY KEY (user_id, token_id)
);
