-- name: ExistsBentoWithNameOwnedByUser :one
SELECT EXISTS(SELECT 1 FROM bentos WHERE name = ? AND user_id = ?);

-- name: NewBento :one
INSERT INTO bentos (user_id, name, created_at, updated_at)
VALUES (?, ?, ?, ?) RETURNING id;

-- name: GetBentoByIDWithPermissions :one
SELECT b.*, p.bytes FROM bentos b
LEFT JOIN bento_permissions p ON p.user_id = ? AND p.bento_id = b.id
WHERE b.id = ?;

-- name: GetBentoWithIDOwnedByUser :one
SELECT * FROM bentos WHERE id = ? AND user_id = ?;

-- name: AddIngredientToBento :exec
INSERT INTO bento_ingredients (bento_id, name, value, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: RemoveIngredientFromBento :execrows
DELETE FROM bento_ingredients WHERE bento_id = ? AND id = ?;

-- name: SetBentoIngredient :exec
INSERT INTO bento_ingredients (bento_id, name, value, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT DO UPDATE SET
    value = excluded.value,
    updated_at = excluded.updated_at;

-- name: GetBentoIngredients :many
SELECT id, name, CAST(value AS TEXT) FROM bento_ingredients
WHERE bento_id = ?;
