-- name: CreateUserToken :one
INSERT INTO user_tokens (
  user_id,
  token_type,
  expires_at,
  created_at
) VALUES (
  ?1, ?2, ?3, ?4
)
RETURNING user_token_id, user_id, token_type, expires_at, created_at;

-- name: GetUserTokenByType :one
SELECT user_token_id, user_id, token_type, expires_at, created_at 
FROM user_tokens
WHERE user_id = ?1 AND token_type = ?2
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteUserToken :exec
DELETE FROM user_tokens
WHERE user_id = ?1 AND token_type = ?2;

-- name: DeleteAllUserTokens :exec
DELETE FROM user_tokens
WHERE user_id = ?1;

-- name: DeleteUserTokensByType :exec
DELETE FROM user_tokens
WHERE user_id = ?1 AND token_type = ?2;