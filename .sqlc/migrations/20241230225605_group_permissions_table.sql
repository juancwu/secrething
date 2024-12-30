-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS group_permissions(
    group_id TEXT NOT NULL,
    bento_id TEXT NOT NULL,
    level INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    CONSTRAINT pk_group_permissions PRIMARY KEY (group_id, bento_id),
    CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS group_permissions;
-- +goose StatementEnd
