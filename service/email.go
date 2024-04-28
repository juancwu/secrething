package service

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/matoous/go-nanoid/v2"
	"github.com/resend/resend-go/v2"

	"github.com/juancwu/konbini/database"
	"github.com/juancwu/konbini/env"
)

type EmailVerification struct {
	Id        int64
	RefId     string
	Status    string // one of "pending" | "opened" | "verified"
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	EMAIL_STATUS_PENDING  = "pending"
	EMAIL_STATUS_OPENED   = "opened"
	EMAIL_STATUS_VERIFIED = "verified"
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

func CreateEmailVerification(userId int64) (string, error) {
	log.Info("Get reference id for email verification")
	refId, err := gonanoid.New(16)
	if err != nil {
		log.Errorf("Error getting reference id for email verification: %v\n", err)
		return "", err
	}

	log.Info("Creating email verification...")
	res, err := database.DB().Exec("INSERT INTO email_verifications (ref_id, user_id) VALUES ($1, $2);", refId, userId)
	if err != nil {
		log.Errorf("Error creating email verification: %v\n", err)
		return "", err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Errorf("Error getting the count for email verification inserted: %v\n", err)
	} else if count > 0 {
		log.Info("Email verification created.", "count", count)
	}

	return refId, nil
}

func GetEmailVerification(refId string) (*EmailVerification, error) {
	log.Info("Get email verification with refId.", "refId", refId)
	row := database.DB().QueryRow("SELECT id, ref_id, status, user_id, created_at, updated_at FROM email_verifications WHERE ref_id = $1;", refId)
	if row.Err() != nil {
		log.Errorf("Error querying email verification: %v\n", row.Err())
		return nil, row.Err()
	}

	log.Info("Scanning email verification values...")
	ev := EmailVerification{}
	err := row.Scan(&ev.Id, &ev.RefId, &ev.Status, &ev.UserId, &ev.CreatedAt, &ev.UpdatedAt)
	if err != nil {
		log.Errorf("Error scanning email verification values: %v\n", err)
		return nil, err
	}

	return &ev, nil
}
