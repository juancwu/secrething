-- name: CreatePermission :one
INSERT INTO permissions (
  permission_id,
  vault_id,
  grantee_type,
  grantee_id,
  permission_bits,
  granted_by,
  created_at,
  updated_at
) VALUES (
  ?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8
)
RETURNING permission_id, vault_id, grantee_type, grantee_id, permission_bits, granted_by, created_at, updated_at;

-- name: GetPermissionByID :one
SELECT permission_id, vault_id, grantee_type, grantee_id, permission_bits, granted_by, created_at, updated_at
FROM permissions
WHERE permission_id = ?1;

-- name: GetPermissionsByVaultID :many
SELECT permission_id, vault_id, grantee_type, grantee_id, permission_bits, granted_by, created_at, updated_at
FROM permissions
WHERE vault_id = ?1;

-- name: GetPermissionByGrantee :one
SELECT permission_id, vault_id, grantee_type, grantee_id, permission_bits, granted_by, created_at, updated_at
FROM permissions
WHERE vault_id = ?1 AND grantee_type = ?2 AND grantee_id = ?3;

-- name: GetPermissionsByGrantee :many
SELECT permission_id, vault_id, grantee_type, grantee_id, permission_bits, granted_by, created_at, updated_at
FROM permissions
WHERE grantee_type = ?1 AND grantee_id = ?2;

-- name: UpdatePermission :one
UPDATE permissions
SET permission_bits = ?2,
    updated_at = ?3
WHERE permission_id = ?1
RETURNING permission_id, vault_id, grantee_type, grantee_id, permission_bits, granted_by, created_at, updated_at;

-- name: DeletePermission :exec
DELETE FROM permissions
WHERE permission_id = ?1;

-- name: DeletePermissionsByVaultID :exec
DELETE FROM permissions
WHERE vault_id = ?1;

-- name: DeletePermissionsByGrantee :exec
DELETE FROM permissions
WHERE grantee_type = ?1 AND grantee_id = ?2;