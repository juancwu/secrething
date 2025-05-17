package db

import "go.jetify.com/typeid"

type VaultIDPrefix struct{}

func (VaultIDPrefix) Prefix() string { return "vault" }

type VaultID struct {
	typeid.TypeID[VaultIDPrefix]
}

type UserIDPrefix struct{}

func (UserIDPrefix) Prefix() string { return "user" }

type UserID struct {
	typeid.TypeID[UserIDPrefix]
}
