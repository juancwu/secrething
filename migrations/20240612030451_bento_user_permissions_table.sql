-- +goose Up
-- +goose StatementBegin
CREATE TABLE bento_user_permissions (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    bento_id UUID NOT NULL,
    permissions SMALLINT NOT NULL DEFAULT 7, -- default to read access
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_bento_user_permission UNIQUE (user_id, bento_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(id) ON DELETE CASCADE
);

CREATE TRIGGER update_bento_user_permissions_updated_at
BEFORE UPDATE ON bento_user_permissions
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_bento_user_permissions_updated_at ON bento_user_permissions;
DROP TABLE bento_user_permissions;
-- +goose StatementEnd
