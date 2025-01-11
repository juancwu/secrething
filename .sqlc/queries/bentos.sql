-- name: ExistsBentoWithNameOwnedByUser :one
SELECT EXISTS(SELECT 1 FROM bentos WHERE name = ? AND user_id = ?);

-- name: NewBento :one
INSERT INTO bentos (user_id, name, created_at, updated_at)
VALUES (?, ?, ?, ?) RETURNING id;

-- name: GetBentoWithIDOwnedByUser :one
SELECT * FROM bentos WHERE id = ? AND user_id = ?;

-- name: AddIngredientToBento :exec
INSERT INTO bento_ingredients (bento_id, name, value, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: SetBentoIngredient :exec
INSERT INTO bento_ingredients (bento_id, name, value, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT DO UPDATE SET
    value = excluded.value,
    updated_at = excluded.updated_at;
