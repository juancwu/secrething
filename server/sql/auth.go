package sql

import (
	_ "embed"
)

//go:embed files/auth/create-user.sql
var CreateUser string

//go:embed files/auth/get-user-with-email.sql
var GetUserWithEmail string
