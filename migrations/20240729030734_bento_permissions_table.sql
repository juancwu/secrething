-- +goose Up
-- +goose StatementBegin
CREATE TABLE bento_permissions (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    bento_id UUID NOT NULL,
    permissions INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_bento_permissions_user_bento UNIQUE (user_id, bento_id),
    CONSTRAINT fk_bento_permission_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_bento_permission_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(id) ON DELETE CASCADE
);
CREATE TRIGGER update_bento_permissions_updated_at
BEFORE UPDATE ON bento_permissions
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column(); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_bento_permissions_updated_at ON bento_permissions;
DROP TABLE bento_permissions;
-- +goose StatementEnd
