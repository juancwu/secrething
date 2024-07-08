package email

import (
	"bytes"
	"context"
	"os"

	"github.com/resend/resend-go/v2"
)

// Send is the main method to use when sending an email from the BE.
func Send(subject, from string, to []string, html string) (*resend.SendEmailResponse, error) {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	params := &resend.SendEmailRequest{
		From:    from,
		To:      to,
		Subject: subject,
		Html:    html,
	}
	return client.Emails.Send(params)
}

// RenderVerifiationEmail renders the html for verification email.
func RenderVerifiationEmail(name, url string) (string, error) {
	var html bytes.Buffer
	err := verificationEmailTempl(name, url).Render(context.Background(), &html)
	if err != nil {
		return "", err
	}
	return html.String(), nil
}
