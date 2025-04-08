package db

import "github.com/sumup/typeid"

type UserIDPrefix struct{}

func (UserIDPrefix) Prefix() string {
	return "user"
}

type UserID = typeid.Sortable[UserIDPrefix]

type TokenIDPrefix struct{}

func (TokenIDPrefix) Prefix() string {
	return "token"
}

type TokenID = typeid.Sortable[TokenIDPrefix]
