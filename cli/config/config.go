package config

import "fmt"

var backendUrl string

func Init() {
	if backendUrl == "" {
		backendUrl = "http://localhost:3000/api/v1"
	}
}

func BackendUrl(path string) string {
	if path == "" {
		return backendUrl
	}

	return fmt.Sprintf("%s/%s", backendUrl, path)
}

type Auth struct {
	Token         string
	TokenType     string
	TOTP          bool
	EmailVerified bool
}

var _auth *Auth = nil

func SetAuth(auth Auth) {
	_auth = &auth
}

func GetAuth() *Auth {
	return _auth
}
