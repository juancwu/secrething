package auth

import (
	"fmt"
	"strings"
)

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

// Is is a utility function that matches if the current AuthServiceError is of the given type.
func (e AuthServiceError) Is(errType string) bool {
	return e.Type == errType
}

const (
	UserAlreadyExistsErr string = "user_already_exists"
)

func NewUserAlreadyExistsErr(email string) AuthServiceError {
	return AuthServiceError{
		Type:    UserAlreadyExistsErr,
		Message: "User with email '{0}' already exists.",
		Params:  []interface{}{email},
	}
}
