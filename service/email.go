package service

import (
	"github.com/charmbracelet/log"
	"github.com/matoous/go-nanoid/v2"
	"github.com/resend/resend-go/v2"

	"github.com/juancwu/konbini/database"
	"github.com/juancwu/konbini/env"
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

func CreateEmailVerification(userId int64) error {
	log.Info("Get reference id for email verification")
	refId, err := gonanoid.New(16)
	if err != nil {
		log.Errorf("Error getting reference id for email verification: %v\n", err)
		return err
	}

	log.Info("Creating email verification...")
	res, err := database.DB().Exec("INSERT INTO email_verifications (ref_id, user_id) VALUES ($1, $2);", refId, userId)
	if err != nil {
		log.Errorf("Error creating email verification: %v\n", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Errorf("Error getting the count for email verification inserted: %v\n", err)
		return nil
	}

	log.Info("Email verification created.", "count", count)

	return nil
}
