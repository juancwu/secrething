-- name: ExistsBentoWithNameOwnedByUser :one
SELECT EXISTS(SELECT 1 FROM bentos WHERE name = ? AND user_id = ?);

-- name: NewBento :one
INSERT INTO bentos (user_id, name, created_at, updated_at)
VALUES (?, ?, ?, ?) RETURNING id;

-- name: AddIngridientToBento :exec
INSERT INTO bento_ingridients (bento_id, name, value, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);
