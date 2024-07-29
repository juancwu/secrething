-- +goose Up
-- +goose StatementBegin
INSERT INTO bento_permissions (user_id, bento_id, permissions)
-- giving read/write/share/grant-share permissions to all owners
SELECT b.owner_id, b.id, 1111
FROM bentos AS b
WHERE NOT EXISTS (
    SELECT 1
    FROM bento_permissions AS perms
    WHERE perms.user_id = b.owner_id
    AND perms.bento_id = b.id
)
ON CONFLICT (user_id, bento_id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'This migration is more to create the missing bento permissions for existing bentos so it wont remove any existing bento permissions.';
-- +goose StatementEnd
