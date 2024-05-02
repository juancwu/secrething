CREATE TABLE lookup_user_personal_bentos (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    bento_id INTEGER NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES personal_bentos(id),

    CONSTRAINT unique_user_bento UNIQUE (user_id, bento_id)
);
