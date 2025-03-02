-- name: NewRecoveryCodes :exec
INSERT INTO totp_recovery_codes
(user_id, created_at, code)
VALUES
(?1, ?2, ?3),
(?1, ?2, ?4),
(?1, ?2, ?5),
(?1, ?2, ?6),
(?1, ?2, ?7),
(?1, ?2, ?8);

-- name: RemoveUserRecoveryCodes :exec
DELETE FROM totp_recovery_codes
WHERE user_id = ?;

-- name: GetRecoveryCode :one
SELECT * FROM totp_recovery_codes
WHERE user_id = ? AND code = ?;

-- name: UseRecoveryCode :exec
UPDATE totp_recovery_codes SET used = true WHERE user_id = ? AND code = ?;
