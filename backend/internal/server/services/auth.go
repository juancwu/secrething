package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/utils"
	"github.com/sumup/typeid"
)

type AuthService struct{}

var authService *AuthService

func NewAuthService() *AuthService {
	if authService == nil {
		authService = &AuthService{}
	}
	return authService
}

func (s *AuthService) CreateUser(ctx context.Context, email, password string, name *string) (*db.User, error) {
	exists, err := s.ExistsUser(ctx, email)
	if err != nil && sql.ErrNoRows != err {
		return nil, err
	}

	if exists {
		return nil, NewUserAlreadyExistsErr(email)
	}

	q, err := db.Query()
	if err != nil {
		return nil, err
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	timestamp := utils.FormatRFC3339NanoFixed(time.Now())

	userID, err := typeid.New[db.UserID]()
	if err != nil {
		return nil, err
	}

	user, err := q.CreateUser(ctx, db.CreateUserParams{
		UserID:       userID,
		Email:        email,
		PasswordHash: hashed,
		Name:         name,
		CreatedAt:    timestamp,
		UpdatedAt:    timestamp,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) ExistsUser(ctx context.Context, email string) (bool, error) {
	q, err := db.Query()
	if err != nil {
		// Default return true to avoid mistakenly proceed with other operations on error
		return true, err
	}
	_, err = q.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return true, nil
}

// AuthenticateUser authenticates a user with the given email and password.
// It returns the user if successful, otherwise an error.
func (s *AuthService) AuthenticateUser(ctx context.Context, email, password string) (*db.User, error) {
	q, err := db.Query()
	if err != nil {
		return nil, err
	}

	// Get the user from the database
	user, err := q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewUserNotFoundErr(email)
		}
		return nil, err
	}

	// Check if the account is locked
	if user.AccountLockedUntil != nil {
		lockedUntil, err := utils.ParseRFC3339NanoStr(*user.AccountLockedUntil)
		if err != nil {
			return nil, err
		}

		if time.Now().Before(lockedUntil) {
			return nil, NewAccountLockedErr(*user.AccountLockedUntil)
		}
	}

	// Verify the password
	isValid, err := utils.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, err
	}

	if !isValid {
		// Increment failed login attempts
		now := time.Now()
		nowStr := utils.FormatRFC3339NanoFixed(now)
		var lockedUntil *string

		// Lock account for 30 minutes after 5 failed attempts
		if user.FailedLoginAttempts != nil && *user.FailedLoginAttempts >= 4 {
			lockedUntilTime := now.Add(30 * time.Minute)
			lockedUntilStr := utils.FormatRFC3339NanoFixed(lockedUntilTime)
			lockedUntil = &lockedUntilStr
		}

		_, err = q.UpdateFailedLoginAttempt(ctx, db.UpdateFailedLoginAttemptParams{
			UserID:             user.UserID,
			LastFailedLoginAt:  &nowStr,
			AccountLockedUntil: lockedUntil,
			UpdatedAt:          nowStr,
		})
		if err != nil {
			return nil, err
		}

		return nil, NewInvalidCredentialsErr()
	}

	// Reset failed login attempts on successful login
	_, err = q.ResetFailedLoginAttempts(ctx, db.ResetFailedLoginAttemptsParams{
		UserID:    user.UserID,
		UpdatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type AuthServiceError struct {
	Type    string
	Message string
	Params  []interface{}
}

// Error implements the error interface
func (e AuthServiceError) Error() string {
	var result string = e.Message
	// Replace positional placeholders with parameter values
	for i, param := range e.Params {
		placeholder := fmt.Sprintf("{%d}", i)
		stringValue := fmt.Sprintf("%v", param)
		result = strings.Replace(result, placeholder, stringValue, -1)
	}
	return result
}

// IsType is a utility function that matches if the current AuthServiceError is of the given type.
func (e AuthServiceError) IsType(errType string) bool {
	return e.Type == errType
}

const (
	AuthServiceErrUserAlreadyExists      string = "auth_service_user_already_exists"
	AuthServiceErrInvalidCredentials     string = "auth_service_invalid_credentials"
	AuthServiceErrUserNotFound           string = "auth_service_user_not_found"
	AuthServiceErrAccountLocked          string = "auth_service_account_locked"
	AuthServiceErrEmailNotVerified       string = "auth_service_email_not_verified"
	AuthServiceErrAuthenticationRequired string = "auth_service_authentication_required"
)

func NewUserAlreadyExistsErr(email string) AuthServiceError {
	return AuthServiceError{
		Type:    AuthServiceErrUserAlreadyExists,
		Message: "User with email '{0}' already exists.",
		Params:  []interface{}{email},
	}
}

func NewInvalidCredentialsErr() AuthServiceError {
	return AuthServiceError{
		Type:    AuthServiceErrInvalidCredentials,
		Message: "Invalid email or password.",
		Params:  []interface{}{},
	}
}

func NewUserNotFoundErr(email string) AuthServiceError {
	return AuthServiceError{
		Type:    AuthServiceErrUserNotFound,
		Message: "User with email '{0}' not found.",
		Params:  []interface{}{email},
	}
}

func NewAccountLockedErr(lockedUntil string) AuthServiceError {
	return AuthServiceError{
		Type:    AuthServiceErrAccountLocked,
		Message: "Account is locked until {0}.",
		Params:  []interface{}{lockedUntil},
	}
}

func NewEmailNotVerifiedErr() AuthServiceError {
	return AuthServiceError{
		Type:    AuthServiceErrEmailNotVerified,
		Message: "Email is not verified. Please check your email for a verification link.",
		Params:  []interface{}{},
	}
}
