package service

import (
	"github.com/charmbracelet/log"
	"github.com/juancwu/konbini/env"
	"github.com/resend/resend-go/v2"
)

func SendEmail(from, to, subject, body string) (string, error) {
	log.Info("Creating resend client...")
	client := resend.NewClient(env.Values().RESEND_API_KEY)

	log.Info("Creating email parameters...")
	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	log.Info("Sending email...", "from", from, "to", to, "subject", subject)
	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Error("Error sending email.", "from", from, "to", to, "subject", subject)
		return "", err
	}
	log.Info("Email sent!", "from", from, "to", to, "subject", subject)

	return sent.Id, nil
}
