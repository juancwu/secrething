package services

import (
	"bytes"
	"context"
	"fmt"
	"konbini/server/config"
	"konbini/server/views"

	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog/log"
)

// SendEmail sends an email via the Resend Client. This is the base function and
// ideally not used directly but instead as the only step where an email is sent.
func SendEmail(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	c, err := config.Global()
	if err != nil {
		return nil, err
	}
	// skip sending emails in testing environment
	if c.IsTesting() {
		return &resend.SendEmailResponse{Id: ""}, nil
	}
	// change the destination email in development to avoid
	// hurting the domain's reputation.
	if c.IsDevelopment() {
		params.To = []string{"delivered@resend.dev"}
	}
	client := resend.NewClient(c.GetResendApiKey())
	sent, err := client.Emails.SendWithContext(ctx, params)
	return sent, err
}

// SendBatchEmails sends a batch of emails at the same time.
// This function should not be used directly but instead when the batch
// is ready to send, use it.
func SendBatchEmails(ctx context.Context, params []*resend.SendEmailRequest) (*resend.BatchEmailResponse, error) {
	cfg, err := config.Global()
	if err != nil {
		return nil, err
	}

	if cfg.IsTesting() {
		return nil, nil
	}

	if cfg.IsDevelopment() {
		// change all the destination emails
		for _, p := range params {
			p.To = []string{"delivered@resend.dev"}
		}
	}

	client := resend.NewClient(cfg.GetResendApiKey())
	return client.Batch.SendWithContext(ctx, params)
}

// SendVerificationEmail sends an email verification for users to verify their email.
func SendVerificationEmail(ctx context.Context, to string, token string) (*resend.SendEmailResponse, error) {
	c, err := config.Global()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/auth/email/verify?token=%s", c.GetBackendUrl(), token)
	component := views.VerificationEmail(url)
	var buffer bytes.Buffer
	err = component.Render(ctx, &buffer)
	if err != nil {
		return nil, err
	}

	params := &resend.SendEmailRequest{
		From:    c.GetVerifyEmailAddress(),
		To:      []string{to},
		Subject: "Verify Your Email",
		Html:    buffer.String(),
		Text: fmt.Sprintf(
			`Thanks for using Konbini!

Please verify your email by opening the following link in a browser:

%s

Do not reply to this email. This email is not monitored.`,
			url,
		),
	}

	res, err := SendEmail(ctx, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type SendGroupInvitationEmailsParams struct {
	InvitorName string
	GroupName   string
	Users       []struct {
		Name  string
		Token string
		Email string
	}
}

func SendGroupInvitationEmails(ctx context.Context, params SendGroupInvitationEmailsParams) (*resend.BatchEmailResponse, error) {
	cfg, err := config.Global()
	if err != nil {
		return nil, err
	}

	batch := make([]*resend.SendEmailRequest, len(params.Users))
	for i, u := range params.Users {
		url := fmt.Sprintf("%s/api/v1/group/invitation/accept?token=%s", cfg.GetBackendUrl(), u.Token)

		var buf bytes.Buffer
		err = views.GroupInvitationEmail(params.InvitorName, u.Name, params.GroupName, url).Render(ctx, &buf)
		if err != nil {
			log.Error().Err(err).Str("to", u.Email).Msg("Failed to render group invitation email.")
			continue
		}

		req := &resend.SendEmailRequest{
			From:    cfg.GetGroupInvitationEmailAddress(),
			To:      []string{u.Email},
			Subject: fmt.Sprintf("Join Group [%s]", params.GroupName),
			Html:    buf.String(),
		}

		batch[i] = req
	}

	return SendBatchEmails(ctx, batch)
}
