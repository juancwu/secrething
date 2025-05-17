-- name: GetVaultByName :one
SELECT
    vault_id,
    vault_name,
    vault_owner_id,
    created_at,
    updated_at
FROM vaults WHERE vault_name = ?1;

-- name: GetVaultByID :one
SELECT
    vault_id,
    vault_name,
    vault_owner_id,
    created_at,
    updated_at
FROM vaults WHERE vault_id = ?1;

-- name: GetAllVaultsByOwner :many
SELECT
    vault_id,
    vault_name,
    vault_owner_id,
    created_at,
    updated_at
FROM vaults WHERE vault_owner_id = ?1;
