package vault

import (
	"errors"

	"github.com/zalando/go-keyring"
)

type Vault struct {
	Token string
}

const (
	keyringService = "konbini"
	keyringUser    = "user"
)

var vault *Vault = nil

func Init() error {
	token, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		return err
	}

	vault = &Vault{
		Token: token,
	}

	return nil
}

// Get gets the current instance of the vault. If not initialized, then nil is returned.
func Get() *Vault {
	return vault
}

// Token retrieves the current auth token. If no token, then it returns an empty string.
func Token() string {
	assertVault()
	return vault.Token
}

// SetToken updates the vault's token value and also the OS keyring
func SetToken(token string) error {
	assertVault()
	err := keyring.Set(keyringService, keyringUser, token)
	if err != nil {
		return err
	}
	vault.Token = token
	return nil
}

// assertVault make sures the vault has been initialized before usage. Panics if not initiliazed
func assertVault() {
	if vault == nil {
		panic(errors.New("Vault used before initialized"))
	}
}
