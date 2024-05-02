CREATE EXTENSION pgcrypto;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) DEFAULT '',
    last_name VARCHAR(255) DEFAULT '',
    email TEXT UNIQUE NOT NULL
);
