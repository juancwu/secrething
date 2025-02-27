package handlers

import (
	"context"
	"errors"
	"github.com/juancwu/konbini/server/db"
	"github.com/juancwu/konbini/server/middlewares"
	"github.com/juancwu/konbini/server/services"
	"github.com/juancwu/konbini/server/utils"
	"time"
)

var ErrUsedRecoveryCode error = errors.New("Recovery code has already been used")

// verifyRecoveryCode validates a recovery code for a user and marks it as used if valid
func verifyRecoveryCode(ctx context.Context, queries *db.Queries, userID, code string) error {
	recoveryCode, err := queries.GetRecoveryCode(ctx, db.GetRecoveryCodeParams{
		UserID: userID,
		Code:   code,
	})
	if err != nil {
		// Record the failed attempt
		middlewares.RecordFailedTOTPAttempt(userID)
		return err
	}

	if recoveryCode.Used {
		// Record the failed attempt with used code
		middlewares.RecordFailedTOTPAttempt(userID)
		return ErrUsedRecoveryCode
	}

	err = queries.UseRecoveryCode(ctx, db.UseRecoveryCodeParams{
		UserID: userID,
		Code:   recoveryCode.Code,
	})
	if err != nil {
		return err
	}

	// Reset the attempt counter on successful verification
	middlewares.ResetTOTPAttempts(userID)

	return nil
}

func newAuthToken(ctx context.Context, queries *db.Queries, userID string, tokType services.TokenType) (*services.AuthToken, error) {
	// If this is a successful auth operation, reset any TOTP attempt counters
	if tokType == services.FULL_USER_TOKEN_TYPE {
		middlewares.ResetTOTPAttempts(userID)
	}

	now := time.Now()
	exp := now.Add(time.Hour * 24 * 7)
	authToken, err := queries.NewAuthToken(ctx, db.NewAuthTokenParams{
		UserID:    userID,
		TokenType: tokType.String(),
		CreatedAt: utils.FormatRFC3339NanoFixed(now),
		ExpiresAt: utils.FormatRFC3339NanoFixed(exp),
	})
	if err != nil {
		return nil, err
	}
	return services.NewAuthToken(authToken.ID, userID, tokType, exp)
}
