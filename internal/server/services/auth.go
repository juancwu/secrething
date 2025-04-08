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

func (*AuthService) CreateUser(ctx context.Context, email, password string, name *string) (*db.User, error) {
	exists, err := ExistsUser(ctx, email)
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

func ExistsUser(ctx context.Context, email string) (bool, error) {
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
	ErrUserAlreadyExists string = "user_already_exists"
)

func NewUserAlreadyExistsErr(email string) AuthServiceError {
	return AuthServiceError{
		Type:    ErrUserAlreadyExists,
		Message: "User with email '{0}' already exists.",
		Params:  []interface{}{email},
	}
}
