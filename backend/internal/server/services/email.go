package services

import (
	"bytes"
	"context"
	"fmt"

	"github.com/juancwu/secrething/internal/server/config"
	"github.com/juancwu/secrething/internal/server/db"
	"github.com/juancwu/secrething/internal/server/templates"
	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client *resend.Client
	cfg    config.EmailConfig
}

var emailService *EmailService

func NewEmailService() *EmailService {
	if emailService == nil {
		cfg := config.Email()
		client := resend.NewClient(cfg.ResendApiKey)
		emailService = &EmailService{
			client: client,
			cfg:    cfg,
		}
	}
	return emailService
}

// SendAccountVerificationEmail sends an email to verify an account.
func (s *EmailService) SendAccountVerificationEmail(ctx context.Context, email string, tokenID db.TokenID) error {
	serverCfg := config.Server()
	// Create verification URL with the token
	verificationURL := fmt.Sprintf("%s/api/auth/account/verify?token=%s", serverCfg.URL, tokenID)

	// Render the template to HTML
	var htmlBuffer bytes.Buffer
	err := templates.AccountVerificationEmail(verificationURL).Render(ctx, &htmlBuffer)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Create and send the email
	params := &resend.SendEmailRequest{
		From:    s.cfg.VerifyAddress,
		To:      []string{email},
		Subject: "Verify Your SecretHing Account",
		Html:    htmlBuffer.String(),
	}

	if serverCfg.Env == config.ServerEnvDevelopment {
		params.To = []string{"delivered@resend.dev"}
	}

	_, err = s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}
