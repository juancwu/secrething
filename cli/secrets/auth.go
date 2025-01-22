package secrets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zalando/go-keyring"

	"konbini/cli/config"
	commonApi "konbini/common/api"
)

const (
	serviceName string = "konbini"
	serviceUser string = "user"
)

type User struct {
	token         string
	tokenType     string
	email         string
	totpSet       bool
	emailVerified bool
}

var user User

// CheckAuth checks if the current auth token is still valid or not.
func CheckAuth() (err error) {
	token, err := keyring.Get(serviceName, serviceUser)
	if err != nil {
		return err
	}

	b, err := json.Marshal(commonApi.CheckAuthTokenRequest{AuthToken: token})
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
	var resBody commonApi.CheckAuthResponse
	json.Unmarshal(data, &resBody)

	// a new token is issued if the user has been active for at least 2 days within the 7 days
	// the old token was issued.

	user = User{
		token:         resBody.AuthToken,
		tokenType:     resBody.TokenType,
		email:         resBody.Email,
		emailVerified: resBody.EmailVerified,
		totpSet:       resBody.TOTP,
	}

	// save new token in keyring
	SaveCredentials(user.token)

	return nil
}

// AuthToken returns the user auth token.
// This token might be an empty string depending on how the CheckAuth success.
func AuthToken() string {
	return user.token
}

// UserEmail returns the user's email.
// This token might be an empty string depending on how the CheckAuth success.
func UserEmail() string {
	return user.email
}

func TokenType() string {
	return user.tokenType
}

func EmailVerified() bool {
	return user.emailVerified
}

func TOTPSet() bool {
	return user.totpSet
}

func SetAuthToken(s string) {
	user.token = s
}

func SetUserEmail(s string) {
	user.email = s
}

func SetTokenType(s string) {
	user.tokenType = s
}

func SetEmailVerified(b bool) {
	user.emailVerified = b
}

func SetTOTP(b bool) {
	user.totpSet = b
}

// SaveCredentials saves the auth token and user email in system keyring.
func SaveCredentials(token string) error {
	return keyring.Set(serviceName, serviceUser, token)
}
