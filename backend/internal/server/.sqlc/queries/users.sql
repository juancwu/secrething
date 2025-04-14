-- name: CreateUser :one
INSERT INTO users (
  user_id,
  email,
  password_hash,
  name,
  created_at,
  updated_at
) VALUES (
  ?1, ?2, ?3, ?4, ?5, ?6
)
RETURNING user_id, email, password_hash, name, email_verified, totp_secret, totp_enabled, 
  failed_login_attempts, last_failed_login_at, account_locked_until, created_at, updated_at;

-- name: GetUserByID :one
SELECT user_id, email, password_hash, name, email_verified, totp_secret, totp_enabled, 
  failed_login_attempts, last_failed_login_at, account_locked_until, created_at, updated_at
FROM users
WHERE user_id = ?1;

-- name: GetUserByEmail :one
SELECT user_id, email, password_hash, name, email_verified, totp_secret, totp_enabled, 
  failed_login_attempts, last_failed_login_at, account_locked_until, created_at, updated_at
FROM users
WHERE email = ?1;

-- name: UpdateUserEmailVerification :one
UPDATE users
SET email_verified = ?2, updated_at = ?3
WHERE user_id = ?1
RETURNING user_id, email, password_hash, name, email_verified, totp_secret, totp_enabled, 
  failed_login_attempts, last_failed_login_at, account_locked_until, created_at, updated_at;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = ?2, updated_at = ?3
WHERE user_id = ?1
RETURNING user_id, email, password_hash, name, email_verified, totp_secret, totp_enabled, 
  failed_login_attempts, last_failed_login_at, account_locked_until, created_at, updated_at;

-- name: UpdateUserName :one
UPDATE users
SET name = ?2, updated_at = ?3
WHERE user_id = ?1
RETURNING user_id, email, password_hash, name, email_verified, totp_secret, totp_enabled, 
  failed_login_attempts, last_failed_login_at, account_locked_until, created_at, updated_at;

-- name: UpdateUserTOTP :one
UPDATE users
SET totp_secret = ?2, totp_enabled = ?3, updated_at = ?4
WHERE user_id = ?1
RETURNING user_id, email, password_hash, name, email_verified, totp_secret, totp_enabled, 
  failed_login_attempts, last_failed_login_at, account_locked_until, created_at, updated_at;

-- name: UpdateFailedLoginAttempt :one
UPDATE users
SET failed_login_attempts = failed_login_attempts + 1, 
    last_failed_login_at = ?2,
    account_locked_until = CASE 
      WHEN failed_login_attempts + 1 >= 5 THEN ?3
      ELSE account_locked_until
    END,
    updated_at = ?4
WHERE user_id = ?1
RETURNING user_id, email, password_hash, name, email_verified, totp_secret, totp_enabled, 
  failed_login_attempts, last_failed_login_at, account_locked_until, created_at, updated_at;

-- name: ResetFailedLoginAttempts :one
UPDATE users
SET failed_login_attempts = 0, 
    last_failed_login_at = NULL,
    account_locked_until = NULL,
    updated_at = ?2
WHERE user_id = ?1
RETURNING user_id, email, password_hash, name, email_verified, totp_secret, totp_enabled, 
  failed_login_attempts, last_failed_login_at, account_locked_until, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = ?1;
