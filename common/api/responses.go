package api

import (
	"encoding/json"
	"io"
)

type SetupTOTPResponse struct {
	URL string `json:"url"`
}

type LockTOTPResponse struct {
	RecoveryCodes []string `json:"recovery_codes"`
	Token         string   `json:"token"`
	Type          string   `json:"type"`
}

type CheckAuthResponse struct {
	AuthToken string `json:"token"`
	TokenType string `json:"type"`
	Email     string `json:"email"`
	// EmailVerified indicates if email has been verified
	EmailVerified bool `json:"email_verified"`
	// TOTP indicates if TOTP has been setup
	TOTP bool `json:"totp"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type RegisterResponse struct {
	AuthToken string `json:"token"`
	TokenType string `json:"type"`
}

type ErrorResponse struct {
	Code      int      `json:"code"`
	Message   string   `json:"message"`
	Errors    []string `json:"errors,omitempty"`
	RequestId string   `json:"request_id"`
}

func ReadErrorResponse(body io.Reader) (*ErrorResponse, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var res ErrorResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
