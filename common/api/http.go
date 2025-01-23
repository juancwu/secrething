package api

import "fmt"

const (
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"

	MimeApplicationJson = "application/json"
)

// Bearer formats a string valid for the Bearer token format
func Bearer(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}
