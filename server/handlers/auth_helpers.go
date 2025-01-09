package handlers

import (
	"context"
	"konbini/server/db"
	"konbini/server/services"
	"konbini/server/utils"
	"time"
)

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
