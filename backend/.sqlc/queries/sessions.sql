-- name: CreateSession :one
INSERT INTO sessions (
    session_id,
    user_id,
    token_hash,
    expires_at,
    created_at,
    last_used_at
) VALUES (
    ?1,
    ?2,
    ?3,
    ?4,
    ?5,
    ?6
) RETURNING *;

-- name: GetSessionByTokenHash :one
SELECT
    s.session_id,
    s.user_id,
    s.token_hash,
    s.expires_at,
    s.created_at,
    s.last_used_at,
    u.user_id,
    u.email,
    u.password_hash,
    u.first_name,
    u.last_name,
    u.created_at as user_created_at,
    u.updated_at as user_updated_at
FROM sessions s
JOIN users u ON s.user_id = u.user_id
WHERE s.token_hash = ?1 AND s.expires_at > strftime('%Y-%m-%dT%H:%M:%SZ', 'now');

-- name: UpdateSessionLastUsed :exec
UPDATE sessions
SET last_used_at = ?2
WHERE session_id = ?1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_id = ?1;

-- name: DeleteSessionByTokenHash :exec
DELETE FROM sessions
WHERE token_hash = ?1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= strftime('%Y-%m-%dT%H:%M:%SZ', 'now');

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = ?1;