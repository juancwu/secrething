package sql

import (
	_ "embed"
)

//go:embed files/email/create-verification-email.sql
var CreateVerificationEmail string

//go:embed files/email/get-verification-email.sql
var GetVerificationEmail string
