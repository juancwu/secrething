-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bento_ingredients (
    ingredient_id TEXT NOT NULL PRIMARY KEY,
    bento_id TEXT NOT NULL CHECK (bento_id != ''),
    name TEXT NOT NULL CHECK (name != ''),
    value BLOB NOT NULL,
    created_at TEXT NOT NULL CHECK (created_at != ''),
    updated_at TEXT NOT NULL CHECK (updated_at != ''),
    CONSTRAINT unique_bento_ingridient_name UNIQUE (bento_id, name),
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(bento_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bento_ingredients;
-- +goose StatementEnd
