package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"konbini/cli/config"
	"konbini/cli/vault"
	"konbini/common/api"
	"net/http"
	"time"
)

var (
	ErrNoAuthToken error = errors.New("No auth token found in keyring")
)

type AuthStatus struct {
	Token     string
	TokenType string
	Email     string
}

func CheckAuth() (*AuthStatus, error) {
	token := vault.Token()
	if token == "" {
		return nil, ErrNoAuthToken
	}

	reqBody := api.CheckAuthTokenRequest{AuthToken: token}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(reqBodyData)

	req, err := http.NewRequest(http.MethodPost, config.BackendUrl(api.UriCheckToken), reader)
	if err != nil {
		return nil, err
	}
	req.Header.Add(api.HeaderContentType, api.MimeApplicationJson)

	c := http.Client{Timeout: time.Second * 10}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := api.ReadErrorResponse(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Check auth error: %s", body.Message)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resBody api.CheckAuthResponse
	err = json.Unmarshal(data, &resBody)
	if err != nil {
		return nil, err
	}

	return &AuthStatus{
		Token:     resBody.AuthToken,
		TokenType: resBody.TokenType,
		Email:     resBody.Email,
	}, nil
}
