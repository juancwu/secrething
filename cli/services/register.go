package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"konbini/cli/config"
	commonAPI "konbini/common/api"
)

func Register(email string, nickname string, password string) (commonAPI.RegisterResponse, error) {
	body, err := json.Marshal(map[string]string{
		"email":    email,
		"nickname": nickname,
		"password": password,
	})
	if err != nil {
		return commonAPI.RegisterResponse{}, err
	}

	reader := bytes.NewReader(body)

	req, err := http.NewRequest(http.MethodPost, config.BackendUrl()+"/auth/register", reader)
	if err != nil {
		return commonAPI.RegisterResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(req)
	if err != nil {
		return commonAPI.RegisterResponse{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return commonAPI.RegisterResponse{}, err
	}

	if res.StatusCode == http.StatusCreated {
		var resBody commonAPI.RegisterResponse
		err = json.Unmarshal(data, &resBody)
		if err != nil {
			return commonAPI.RegisterResponse{}, err
		}

		return resBody, nil
	}

	var resBody commonAPI.ErrorResponse
	err = json.Unmarshal(data, &resBody)
	if err != nil {
		return commonAPI.RegisterResponse{}, err
	}

	return commonAPI.RegisterResponse{}, fmt.Errorf("Registration Error: %s", resBody.Message)

}
