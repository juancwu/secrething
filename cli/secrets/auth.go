package secrets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"konbini/cli/config"
	"net/http"
	"strings"
	"time"

	"github.com/zalando/go-keyring"
)

const (
	serviceName string = "konbini"
	serviceUser string = "user"
)

var authToken string
var userEmail string
var tokenType string

// CheckAuth checks if the current auth token is still valid or not.
func CheckAuth() (err error) {
	var key string
	key, err = keyring.Get(serviceName, serviceUser)
	if err != nil {
		return err
	}

	parts := strings.Split(key, " ")
	if len(parts) != 2 {
		return fmt.Errorf("Invalid key format")
	}

	b, err := json.Marshal(map[string]string{"auth_token": authToken})
	if err != nil {
		return err
	}

	reader := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPost, config.BackendUrl()+"/auth/token/check", reader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Auth token is not valid")
		return err
	}

	// check if a new token has been issued
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var resBody map[string]string
	json.Unmarshal(data, &resBody)

	// a new token is issued if the user has been active for at least 2 days within the 7 days
	// the old token was issued.
	if newToken, ok := resBody["auth_token"]; ok {
		parts[0] = newToken
	}

	authToken = parts[0]
	userEmail = parts[1]
	tokenType = resBody["type"]

	// save new token in keyring
	SaveCredentials(authToken, userEmail)

	return nil
}

// AuthToken returns the user auth token.
// This token might be an empty string depending on how the CheckAuth success.
func AuthToken() string {
	return authToken
}

// UserEmail returns the user's email.
// This token might be an empty string depending on how the CheckAuth success.
func UserEmail() string {
	return userEmail
}

func TokenType() string {
	return tokenType
}

func SetAuthToken(s string) {
	authToken = s
}

func SetUserEmail(s string) {
	userEmail = s
}

func SetTokenType(s string) {
	tokenType = s
}

// SaveCredentials saves the auth token and user email in system keyring.
func SaveCredentials(token string, email string) error {
	return keyring.Set(serviceName, serviceUser, fmt.Sprintf("%s %s", token, email))
}
