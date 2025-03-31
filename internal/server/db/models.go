// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

type AccessLog struct {
	AccessLogID  string      `db:"access_log_id" json:"access_log_id"`
	UserID       *string     `db:"user_id" json:"user_id"`
	BentoID      *string     `db:"bento_id" json:"bento_id"`
	GroupID      *string     `db:"group_id" json:"group_id"`
	BentoTokenID *string     `db:"bento_token_id" json:"bento_token_id"`
	Action       string      `db:"action" json:"action"`
	Details      interface{} `db:"details" json:"details"`
	AccessedAt   string      `db:"accessed_at" json:"accessed_at"`
}

type Device struct {
	DeviceID            string      `db:"device_id" json:"device_id"`
	UserID              string      `db:"user_id" json:"user_id"`
	DeviceName          interface{} `db:"device_name" json:"device_name"`
	DeviceType          interface{} `db:"device_type" json:"device_type"`
	OsType              interface{} `db:"os_type" json:"os_type"`
	OsVersion           interface{} `db:"os_version" json:"os_version"`
	AppVersion          interface{} `db:"app_version" json:"app_version"`
	BrowserType         interface{} `db:"browser_type" json:"browser_type"`
	BrowserVersion      interface{} `db:"browser_version" json:"browser_version"`
	IpAddress           interface{} `db:"ip_address" json:"ip_address"`
	UserAgent           interface{} `db:"user_agent" json:"user_agent"`
	DeviceFingerprint   interface{} `db:"device_fingerprint" json:"device_fingerprint"`
	IsTrusted           *bool       `db:"is_trusted" json:"is_trusted"`
	FirstSeenAt         interface{} `db:"first_seen_at" json:"first_seen_at"`
	LastActiveAt        interface{} `db:"last_active_at" json:"last_active_at"`
	CreatedAt           string      `db:"created_at" json:"created_at"`
	UpdatedAt           string      `db:"updated_at" json:"updated_at"`
	LastLatitude        interface{} `db:"last_latitude" json:"last_latitude"`
	LastLongitude       interface{} `db:"last_longitude" json:"last_longitude"`
	LastLocationCountry interface{} `db:"last_location_country" json:"last_location_country"`
	LastLocationCity    interface{} `db:"last_location_city" json:"last_location_city"`
}

type FailedLoginAttempt struct {
	AttemptID          string      `db:"attempt_id" json:"attempt_id"`
	UserID             interface{} `db:"user_id" json:"user_id"`
	Email              string      `db:"email" json:"email"`
	IpAddress          string      `db:"ip_address" json:"ip_address"`
	UserAgent          interface{} `db:"user_agent" json:"user_agent"`
	DeviceFingerprint  interface{} `db:"device_fingerprint" json:"device_fingerprint"`
	GeolocationCountry interface{} `db:"geolocation_country" json:"geolocation_country"`
	GeolocationCity    interface{} `db:"geolocation_city" json:"geolocation_city"`
	AttemptTime        string      `db:"attempt_time" json:"attempt_time"`
	FailureReason      string      `db:"failure_reason" json:"failure_reason"`
	DeviceID           interface{} `db:"device_id" json:"device_id"`
}

type SecurityEvent struct {
	EventID           string      `db:"event_id" json:"event_id"`
	UserID            interface{} `db:"user_id" json:"user_id"`
	DeviceID          interface{} `db:"device_id" json:"device_id"`
	TokenID           interface{} `db:"token_id" json:"token_id"`
	IpAddress         interface{} `db:"ip_address" json:"ip_address"`
	DeviceFingerprint interface{} `db:"device_fingerprint" json:"device_fingerprint"`
	EventType         string      `db:"event_type" json:"event_type"`
	Severity          string      `db:"severity" json:"severity"`
	DetailsJson       string      `db:"details_json" json:"details_json"`
	IsResolved        bool        `db:"is_resolved" json:"is_resolved"`
	AlertSent         bool        `db:"alert_sent" json:"alert_sent"`
	AlertSentAt       interface{} `db:"alert_sent_at" json:"alert_sent_at"`
	ResolutionNotes   interface{} `db:"resolution_notes" json:"resolution_notes"`
	CreatedAt         string      `db:"created_at" json:"created_at"`
	ResolvedAt        interface{} `db:"resolved_at" json:"resolved_at"`
}

type SecurityNotification struct {
	NotificationID      string      `db:"notification_id" json:"notification_id"`
	UserID              string      `db:"user_id" json:"user_id"`
	SecurityEventID     interface{} `db:"security_event_id" json:"security_event_id"`
	NotificationType    string      `db:"notification_type" json:"notification_type"`
	NotificationChannel string      `db:"notification_channel" json:"notification_channel"`
	Recipient           string      `db:"recipient" json:"recipient"`
	Content             string      `db:"content" json:"content"`
	IsSent              bool        `db:"is_sent" json:"is_sent"`
	SentAt              interface{} `db:"sent_at" json:"sent_at"`
	ErrorMessage        interface{} `db:"error_message" json:"error_message"`
}

type SecurityNotificationSetting struct {
	UserID                  string      `db:"user_id" json:"user_id"`
	NotifyOnNewDevice       bool        `db:"notify_on_new_device" json:"notify_on_new_device"`
	NotifyOnSuspiciousLogin bool        `db:"notify_on_suspicious_login" json:"notify_on_suspicious_login"`
	NotifyOnFailedAttempts  int64       `db:"notify_on_failed_attempts" json:"notify_on_failed_attempts"`
	FailedAttemptsThreshold int64       `db:"failed_attempts_threshold" json:"failed_attempts_threshold"`
	NotifyOnPasswordChange  bool        `db:"notify_on_password_change" json:"notify_on_password_change"`
	NotifyOnEmailChange     bool        `db:"notify_on_email_change" json:"notify_on_email_change"`
	NotifyOnTotpChange      bool        `db:"notify_on_totp_change" json:"notify_on_totp_change"`
	NotificationEmail       interface{} `db:"notification_email" json:"notification_email"`
}

type Team struct {
	TeamID    string `db:"team_id" json:"team_id"`
	Name      string `db:"name" json:"name"`
	OwnerID   string `db:"owner_id" json:"owner_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type TeamInvitation struct {
	TeamInvitationID string `db:"team_invitation_id" json:"team_invitation_id"`
	UserID           string `db:"user_id" json:"user_id"`
	TeamID           string `db:"team_id" json:"team_id"`
	ResendEmailID    string `db:"resend_email_id" json:"resend_email_id"`
	CreatedAt        string `db:"created_at" json:"created_at"`
	ExpiresAt        string `db:"expires_at" json:"expires_at"`
}

type TeamPermission struct {
	TeamID      string `db:"team_id" json:"team_id"`
	VaultID     string `db:"vault_id" json:"vault_id"`
	Permissions int64  `db:"permissions" json:"permissions"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}

type Token struct {
	TokenID               string      `db:"token_id" json:"token_id"`
	UserID                string      `db:"user_id" json:"user_id"`
	DeviceID              interface{} `db:"device_id" json:"device_id"`
	TokenHash             string      `db:"token_hash" json:"token_hash"`
	TokenType             string      `db:"token_type" json:"token_type"`
	AccessLevel           string      `db:"access_level" json:"access_level"`
	IsActive              bool        `db:"is_active" json:"is_active"`
	IssuedAt              string      `db:"issued_at" json:"issued_at"`
	ExpiresAt             string      `db:"expires_at" json:"expires_at"`
	LastUsedAt            string      `db:"last_used_at" json:"last_used_at"`
	RevokedAt             interface{} `db:"revoked_at" json:"revoked_at"`
	RevocationReason      interface{} `db:"revocation_reason" json:"revocation_reason"`
	ParentTokenID         interface{} `db:"parent_token_id" json:"parent_token_id"`
	IssuedIp              interface{} `db:"issued_ip" json:"issued_ip"`
	IssuedLocationCountry interface{} `db:"issued_location_country" json:"issued_location_country"`
	IssuedLocationCity    interface{} `db:"issued_location_city" json:"issued_location_city"`
	SessionID             interface{} `db:"session_id" json:"session_id"`
	Scope                 interface{} `db:"scope" json:"scope"`
}

type TokenEvent struct {
	EventID          string      `db:"event_id" json:"event_id"`
	TokenID          string      `db:"token_id" json:"token_id"`
	UserID           string      `db:"user_id" json:"user_id"`
	DeviceID         interface{} `db:"device_id" json:"device_id"`
	EventType        string      `db:"event_type" json:"event_type"`
	IpAddress        interface{} `db:"ip_address" json:"ip_address"`
	UserAgent        interface{} `db:"user_agent" json:"user_agent"`
	EventDetailsJson interface{} `db:"event_details_json" json:"event_details_json"`
	CreatedAt        string      `db:"created_at" json:"created_at"`
}

type TotpRecoveryCode struct {
	TotpRecoveryCodeID string  `db:"totp_recovery_code_id" json:"totp_recovery_code_id"`
	UserID             string  `db:"user_id" json:"user_id"`
	Code               string  `db:"code" json:"code"`
	Used               bool    `db:"used" json:"used"`
	CreatedAt          string  `db:"created_at" json:"created_at"`
	UsedAt             *string `db:"used_at" json:"used_at"`
}

type User struct {
	UserID              string  `db:"user_id" json:"user_id"`
	Email               string  `db:"email" json:"email"`
	PasswordHash        string  `db:"password_hash" json:"password_hash"`
	Name                *string `db:"name" json:"name"`
	EmailVerified       bool    `db:"email_verified" json:"email_verified"`
	TotpSecret          *string `db:"totp_secret" json:"totp_secret"`
	TotpEnabled         bool    `db:"totp_enabled" json:"totp_enabled"`
	AccountStatus       string  `db:"account_status" json:"account_status"`
	FailedLoginAttempts *int64  `db:"failed_login_attempts" json:"failed_login_attempts"`
	LastFailedLoginAt   *string `db:"last_failed_login_at" json:"last_failed_login_at"`
	AccountLockedUntil  *string `db:"account_locked_until" json:"account_locked_until"`
	CreatedAt           string  `db:"created_at" json:"created_at"`
	UpdatedAt           string  `db:"updated_at" json:"updated_at"`
}

type UserSetting struct {
	UserSettingID             string `db:"user_setting_id" json:"user_setting_id"`
	UserID                    string `db:"user_id" json:"user_id"`
	AlertOnSuspiciousActivity bool   `db:"alert_on_suspicious_activity" json:"alert_on_suspicious_activity"`
	CreatedAt                 string `db:"created_at" json:"created_at"`
	UpdatedAt                 string `db:"updated_at" json:"updated_at"`
}

type UserVerification struct {
	VerificationID      string      `db:"verification_id" json:"verification_id"`
	UserID              string      `db:"user_id" json:"user_id"`
	VerificationType    string      `db:"verification_type" json:"verification_type"`
	VerificationTokenID interface{} `db:"verification_token_id" json:"verification_token_id"`
	VerificationCode    interface{} `db:"verification_code" json:"verification_code"`
	IsVerified          bool        `db:"is_verified" json:"is_verified"`
	Attempts            int64       `db:"attempts" json:"attempts"`
	ExpiresAt           string      `db:"expires_at" json:"expires_at"`
	VerifiedAt          interface{} `db:"verified_at" json:"verified_at"`
	CreatedAt           string      `db:"created_at" json:"created_at"`
}

type UsersTeam struct {
	UserID    string `db:"user_id" json:"user_id"`
	TeamID    string `db:"team_id" json:"team_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type Vault struct {
	VaultID   string `db:"vault_id" json:"vault_id"`
	UserID    string `db:"user_id" json:"user_id"`
	Name      string `db:"name" json:"name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type VaultIngredient struct {
	IngredientID string `db:"ingredient_id" json:"ingredient_id"`
	VaultID      string `db:"vault_id" json:"vault_id"`
	Name         string `db:"name" json:"name"`
	Value        []byte `db:"value" json:"value"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	UpdatedAt    string `db:"updated_at" json:"updated_at"`
}

type VaultPermission struct {
	UserID      string `db:"user_id" json:"user_id"`
	VaultID     string `db:"vault_id" json:"vault_id"`
	Permissions int64  `db:"permissions" json:"permissions"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}

type VaultToken struct {
	VaultTokenID string  `db:"vault_token_id" json:"vault_token_id"`
	VaultID      string  `db:"vault_id" json:"vault_id"`
	TokenSalt    []byte  `db:"token_salt" json:"token_salt"`
	CreatedBy    string  `db:"created_by" json:"created_by"`
	CreatedAt    string  `db:"created_at" json:"created_at"`
	LastUsedAt   *string `db:"last_used_at" json:"last_used_at"`
	ExpiresAt    *string `db:"expires_at" json:"expires_at"`
}
