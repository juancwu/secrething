-- name: CreateVault :one
INSERT INTO vaults (
  vault_id,
  name,
  description,
  created_by_user_id,
  owner_type,
  owner_id,
  created_at,
  updated_at
) VALUES (
  ?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8
)
RETURNING vault_id, name, description, created_by_user_id, owner_type, owner_id, created_at, updated_at;

-- name: GetVaultByID :one
SELECT vault_id, name, description, created_by_user_id, owner_type, owner_id, created_at, updated_at
FROM vaults
WHERE vault_id = ?1;

-- name: GetVaultsByOwner :many
SELECT vault_id, name, description, created_by_user_id, owner_type, owner_id, created_at, updated_at
FROM vaults
WHERE owner_type = ?1 AND owner_id = ?2;

-- name: GetVaultsByUserPermission :many
SELECT v.vault_id, v.name, v.description, v.created_by_user_id, v.owner_type, v.owner_id, v.created_at, v.updated_at
FROM vaults v
JOIN permissions p ON v.vault_id = p.vault_id
WHERE p.grantee_type = 'user' AND p.grantee_id = ?1;

-- name: GetVaultsByTeamPermission :many
SELECT v.vault_id, v.name, v.description, v.created_by_user_id, v.owner_type, v.owner_id, v.created_at, v.updated_at
FROM vaults v
JOIN permissions p ON v.vault_id = p.vault_id
WHERE p.grantee_type = 'team' AND p.grantee_id = ?1;

-- name: UpdateVault :one
UPDATE vaults
SET name = ?2,
    description = ?3,
    updated_at = ?4
WHERE vault_id = ?1
RETURNING vault_id, name, description, created_by_user_id, owner_type, owner_id, created_at, updated_at;

-- name: DeleteVault :exec
DELETE FROM vaults
WHERE vault_id = ?1;