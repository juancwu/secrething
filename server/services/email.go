package services

// SendEmail sends an email via the Resend Client. This is the base function and
// ideally not used directly but instead as the only step where an email is sent.
func SendEmail(subject string, from string, to []string, body interface{}) error {
	return nil
}

// SendVerificationEmail sends an email verification for users to verify their email.
func SendVerificationEmail(to string, token string) error {
	return nil
}
