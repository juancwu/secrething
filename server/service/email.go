package service

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/matoous/go-nanoid/v2"
	"github.com/resend/resend-go/v2"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/env"
)

const (
	EMAIL_STATUS_PENDING  = "PENDING"
	EMAIL_STATUS_SENT     = "SENT"
	EMAIL_STATUS_OPENED   = "OPENED"
	EMAIL_STATUS_FAILED   = "FAILED"
	EMAIL_STATUS_VERIFIED = "VERIFIED"
)

type EmailVerificationStatus string

type EmailVerification struct {
	Id             int64
	VerificationId string
	Status         EmailVerificationStatus // one of email status constants
	UserId         int64
	ResendEmailId  *string
	EmailSentAT    *time.Time
	ExpiresAt      time.Time
	VerifiedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

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
	verificationId, err := gonanoid.New(16)
	if err != nil {
		log.Errorf("Error getting reference id for email verification: %v\n", err)
		return "", err
	}

	log.Info("Creating email verification...")
	// 24 hours from creation
	expTime := time.Now().In(time.UTC).Add(time.Hour * 24)
	res, err := database.DB().
		Exec("INSERT INTO email_verifications (verification_id, user_id, expires_at) VALUES ($1, $2, $3);", verificationId, userId, expTime)
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

	return verificationId, nil
}

func GetEmailVerification(refId string) (*EmailVerification, error) {
	log.Info("Get email verification with refId.", "refId", refId)
	row := database.DB().QueryRow(
		"SELECT id, verification_id, user_id, resend_email_id, status, email_sent_at, expires_at, verified_at, created_at, updated_at FROM email_verifications WHERE verification_id = $1;",
		refId)
	if row.Err() != nil {
		log.Errorf("Error querying email verification: %v\n", row.Err())
		return nil, row.Err()
	}

	log.Info("Scanning email verification values...")
	ev := EmailVerification{}
	err := row.Scan(
		&ev.Id,
		&ev.VerificationId,
		&ev.UserId,
		&ev.ResendEmailId,
		&ev.Status,
		&ev.EmailSentAT,
		&ev.ExpiresAt,
		&ev.VerifiedAt,
		&ev.CreatedAt,
		&ev.UpdatedAt,
	)
	if err != nil {
		log.Errorf("Error scanning email verification values: %v\n", err)
		return nil, err
	}

	return &ev, nil
}
