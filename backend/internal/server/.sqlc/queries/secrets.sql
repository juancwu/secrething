-- name: CreateSecret :one
INSERT INTO secrets (
  secret_id,
  vault_id,
  name,
  value,
  created_by_user_id,
  created_at,
  updated_at
) VALUES (
  ?1, ?2, ?3, ?4, ?5, ?6, ?7
)
RETURNING secret_id, vault_id, name, value, created_by_user_id, created_at, updated_at;

-- name: GetSecretByID :one
SELECT secret_id, vault_id, name, value, created_by_user_id, created_at, updated_at
FROM secrets
WHERE secret_id = ?1;

-- name: GetSecretsByVaultID :many
SELECT secret_id, vault_id, name, value, created_by_user_id, created_at, updated_at
FROM secrets
WHERE vault_id = ?1;

-- name: UpdateSecret :one
UPDATE secrets
SET name = ?2,
    value = ?3,
    updated_at = ?4
WHERE secret_id = ?1
RETURNING secret_id, vault_id, name, value, created_by_user_id, created_at, updated_at;

-- name: DeleteSecret :exec
DELETE FROM secrets
WHERE secret_id = ?1;

-- name: DeleteSecretsByVaultID :exec
DELETE FROM secrets
WHERE vault_id = ?1;