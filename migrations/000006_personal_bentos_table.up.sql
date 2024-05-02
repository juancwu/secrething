-- personal bentos
CREATE TABLE personal_bentos (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL,
    owner_id INTEGER,
    content BYTEA NOT NULL, -- encrypted
    pub_key BYTEA NOT NULL, -- encrypted

    CONSTRAINT unique_bento_owner UNIQUE INDEX (name, owner_id),
    CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES users(id)
);
