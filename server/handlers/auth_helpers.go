package handlers

import (
	"context"
	"errors"
	"konbini/server/db"
	"konbini/server/services"
	"konbini/server/utils"
	"time"
)

var ErrUsedRecoveryCode error = errors.New("Recovery code has already been used")

func newAuthToken(ctx context.Context, queries *db.Queries, userID string, tokType services.TokenType) (*services.AuthToken, error) {
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

// verifyRecoveryCode is a helper function that verifies if the given recovery code is valid for the user to use.
func verifyRecoveryCode(ctx context.Context, q *db.Queries, userID string, recoveryCode string) error {
	row, err := q.GetRecoveryCode(ctx, db.GetRecoveryCodeParams{
		UserID: userID,
		Code:   recoveryCode,
	})
	if err != nil {
		return err
	}
	if row.Used {
		return ErrUsedRecoveryCode
	}
	err = q.UseRecoveryCode(ctx, db.UseRecoveryCodeParams{
		UserID: userID,
		Code:   recoveryCode,
	})
	if err != nil {
		return err
	}
	return nil
}
