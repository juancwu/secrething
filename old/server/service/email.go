package service

import (
	"database/sql"
	"time"

	"github.com/matoous/go-nanoid/v2"
	"github.com/resend/resend-go/v2"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/env"
	"github.com/juancwu/konbini/server/utils"
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
	Id             string
	VerificationId string
	Status         EmailVerificationStatus // one of email status constants
	UserId         string
	ResendEmailId  *string
	EmailSentAT    *time.Time
	ExpiresAt      time.Time
	VerifiedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func SendEmail(from, to, subject, body string) (string, error) {
	utils.Logger().Info("Creating resend client...")
	client := resend.NewClient(env.Values().RESEND_API_KEY)

	utils.Logger().Info("Creating email parameters...")
	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	utils.Logger().Info("Sending email...", "from", from, "to", to, "subject", subject)
	sent, err := client.Emails.Send(params)
	if err != nil {
		utils.Logger().Error("Error sending email.", "from", from, "to", to, "subject", subject)
		return "", err
	}
	utils.Logger().Info("Email sent!", "from", from, "to", to, "subject", subject)

	return sent.Id, nil
}

func CreateEmailVerification(userId string, tx *sql.Tx) (string, error) {
	utils.Logger().Info("Get reference id for email verification")
	verificationId, err := gonanoid.New(16)
	if err != nil {
		utils.Logger().Errorf("Error getting reference id for email verification: %v\n", err)
		return "", err
	}

	utils.Logger().Info("Creating email verification...")
	// 24 hours from creation
	expTime := time.Now().In(time.UTC).Add(time.Hour * 24)
	res, err := tx.
		Exec("INSERT INTO email_verifications (verification_id, user_id, expires_at) VALUES ($1, $2, $3);", verificationId, userId, expTime)
	if err != nil {
		utils.Logger().Errorf("Error creating email verification: %v\n", err)
		return "", err
	}

	count, err := res.RowsAffected()
	if err != nil {
		utils.Logger().Errorf("Error getting the count for email verification inserted: %v\n", err)
	} else if count > 0 {
		utils.Logger().Info("Email verification created.", "count", count)
	}

	return verificationId, nil
}

func GetEmailVerification(refId string) (*EmailVerification, error) {
	utils.Logger().Info("Get email verification with refId.", "refId", refId)
	row := database.DB().QueryRow(
		"SELECT id, verification_id, user_id, resend_email_id, status, email_sent_at, expires_at, verified_at, created_at, updated_at FROM email_verifications WHERE verification_id = $1;",
		refId)
	if row.Err() != nil {
		utils.Logger().Errorf("Error querying email verification: %v\n", row.Err())
		return nil, row.Err()
	}

	utils.Logger().Info("Scanning email verification values...")
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
		utils.Logger().Errorf("Error scanning email verification values: %v\n", err)
		return nil, err
	}

	return &ev, nil
}
