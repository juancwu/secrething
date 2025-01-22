package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"konbini/cli/config"
	"net/http"
	"time"
)

type RegisterResponse struct {
	Token string
	Type  string
}

func Register(email string, nickname string, password string) (RegisterResponse, error) {
	body, err := json.Marshal(map[string]string{
		"email":    email,
		"nickname": nickname,
		"password": password,
	})
	if err != nil {
		return RegisterResponse{}, err
	}

	reader := bytes.NewReader(body)

	req, err := http.NewRequest(http.MethodPost, config.BackendUrl()+"/auth/register", reader)
	if err != nil {
		return RegisterResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(req)
	if err != nil {
		return RegisterResponse{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return RegisterResponse{}, err
	}
	var resBody map[string]any
	err = json.Unmarshal(data, &resBody)
	if err != nil {
		return RegisterResponse{}, err
	}

	if res.StatusCode != http.StatusCreated {
		return RegisterResponse{}, fmt.Errorf("Registration Error: %s", resBody["message"])
	}

	return RegisterResponse{
		Token: resBody["token"].(string),
		Type:  resBody["type"].(string),
	}, nil
}
