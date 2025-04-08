-- name: CreateToken :one
INSERT INTO tokens (
    token_id,
    user_id,
    token_type,
    client_type,
    expires_at,
    created_at
) VALUES (
    ?1, ?2, ?3, ?4, ?5, ?6
)
RETURNING token_id, user_id, token_type, client_type, expires_at, created_at;

-- name: GetTokenByType :one
SELECT token_id, user_id, token_type, client_type, expires_at, created_at 
FROM tokens
WHERE user_id = ?1 AND token_type = ?2
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE user_id = ?1 AND token_type = ?2;

-- name: DeleteAllTokens :exec
DELETE FROM tokens
WHERE user_id = ?1;

-- name: DeleteTokensByType :exec
DELETE FROM tokens
WHERE user_id = ?1 AND token_type = ?2;
