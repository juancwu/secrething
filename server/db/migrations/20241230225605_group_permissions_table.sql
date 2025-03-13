-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS group_permissions(
    group_id TEXT NOT NULL CHECK (group_id != ''),
    bento_id TEXT NOT NULL CHECK (bento_id != ''),
    permissions TEXT NOT NULL CHECK (permissions != ''),
    created_at TEXT NOT NULL CHECK (created_at != ''),
    updated_at TEXT NOT NULL CHECK (updated_at != ''),
    CONSTRAINT pk_group_permissions PRIMARY KEY (group_id, bento_id),
    CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE,
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(bento_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS group_permissions;
-- +goose StatementEnd
