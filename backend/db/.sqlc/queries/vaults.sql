-- name: GetVault :one
SELECT * FROM vaults WHERE vault_name = ?1 AND owner_id = ?2;
