package handlers

import (
	"context"
	"konbini/server/db"
	"konbini/server/services"
	"time"
)

func newAuthToken(ctx context.Context, queries *db.Queries, userID string, tokType services.TokenType) (*services.AuthToken, error) {
	now := time.Now().UTC()
	exp := now.Add(time.Hour * 24 * 7)
	authToken, err := queries.NewAuthToken(ctx, db.NewAuthTokenParams{
		UserID:    userID,
		TokenType: tokType.String(),
		CreatedAt: now.Format(time.RFC3339Nano),
		ExpiresAt: exp.Format(time.RFC3339Nano),
	})
	if err != nil {
		return nil, err
	}
	return services.NewAuthToken(authToken.ID, userID, tokType, exp)
}
