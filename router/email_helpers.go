package router

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/juancwu/konbini/store"
	"github.com/juancwu/konbini/utils"
	"github.com/juancwu/konbini/views"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

// sendVerificationEmail is a helper function that sends a verification email.
func sendVerificationEmail(email string, firstName string, userId string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// generate code
	code, err := gonanoid.Generate(store.EMAIL_VERIFICATION_CODE_CHR_POOL, store.EMAIL_VERIFICATION_CODE_LEN)
	if err != nil {
		logger.Error("Failed to generate email verification code on new user created.", zap.Error(err))
		return
	}

	// try to send email first
	var html bytes.Buffer
	err = views.VerifyEmail(firstName, fmt.Sprintf("%s/api/v1/account/verify-email?code=%s", os.Getenv("SERVER_URL"), code)).Render(context.Background(), &html)
	if err != nil {
		logger.Error("Failed to render email verification view on new user created.", zap.Error(err))
		return
	}

	// save the email verification in the database
	_, err = store.CreateEmailVerification(code, userId)
	if err != nil {
		logger.Error("Failed to save email verification in database on new user created.", zap.Error(err))
		return
	}

	// send email
	_, err = utils.SendEmail(os.Getenv("NOREPLY_EMAIL"), []string{email}, "[Konbini] Verify Your Email", html.String())
	if err != nil {
		logger.Error("Failed to send email verification on new user created.", zap.Error(err))
		return
	}
}

// sendPasswordResetEmail is a helper function that sends a password reset email.
func sendPasswordResetEmail(email, firstName, resetCode string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// try to send email first
	var html bytes.Buffer
	err := views.ResetPasswordEmail(firstName, resetCode).Render(context.Background(), &html)
	if err != nil {
		logger.Error("Failed to render email verification view on new user created.", zap.Error(err))
		return
	}

	_, err = utils.SendEmail(os.Getenv("NOREPLY_EMAIL"), []string{email}, "[Konbini] Reset Password", html.String())
	if err != nil {
		logger.Error("Failed to send reset password email.", zap.Error(err), zap.String("email", email))
		return
	}
}
