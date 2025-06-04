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

type SessionIDPrefix struct{}

func (SessionIDPrefix) Prefix() string { return "session" }

type SessionID struct {
	typeid.TypeID[SessionIDPrefix]
}
