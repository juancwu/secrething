package utils

import (
	"os"

	"github.com/resend/resend-go/v2"
)

// SendEmail is a utility function that sends an email using resend api.
// The body has to be a valid HTML. Returns the email id that was sent.
func SendEmail(from string, to []string, subject, body string) (string, error) {
	c := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	params := &resend.SendEmailRequest{
		From:    from,
		To:      to,
		Subject: subject,
		Html:    body,
	}
	sent, err := c.Emails.Send(params)
	if err != nil {
		return "", err
	}
	return sent.Id, nil
}
