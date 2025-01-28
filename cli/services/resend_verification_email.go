package services

import (
	"errors"
	"io"
	"konbini/cli/config"
	"konbini/common/api"
	"net/http"
	"time"
)

var (
	ErrMissingAuth error = errors.New("Missing auth object")
)

func ResendVerificationEmail() error {
	auth := config.GetAuth()
	if auth == nil || auth.Token == "" {
		return ErrMissingAuth
	}

	req, err := http.NewRequest(http.MethodPost, config.BackendUrl(api.UriResendVerificationEmail), nil)
	if err != nil {
		return err
	}
	req.Header.Add(api.HeaderAuthorization, api.Bearer(auth.Token))

	c := http.Client{Timeout: time.Second * 30}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(data))
	}

	return nil
}
