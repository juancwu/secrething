package templates

import _ "embed"

//go:embed .static/internal_error.html
var StaticInternalError string

//go:embed .static/verification_error.html
var StaticVerificationError string

//go:embed .static/verification_success.html
var StaticVerificationSuccess string
