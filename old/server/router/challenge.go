package router

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/service"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

func SetupChallengeRoutes(e RouteGroup) {
	e.GET("/challenge", handleGetChallenge, JwtAuthMiddleware)
}

func handleGetChallenge(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	claims, ok := c.Get("claims").(*service.JwtCustomClaims)
	if !ok {
		logger.Error("Failed to cast claims to service.JwtCustomClaims")
		return writeApiError(c, http.StatusInternalServerError, "internal server error")
	}

	// create a new challenge
	randomBytes := make([]byte, 128)
	_, err := rand.Read(randomBytes)
	if err != nil {
		logger.Error("Failed to create random bytes for challenge", zap.Error(err))
		return writeApiError(c, http.StatusInternalServerError, "internal server error")
	}

	hash := sha256.New()
	_, err = hash.Write(randomBytes)
	if err != nil {
		logger.Error("Failed to write random bytes to hash", zap.Error(err))
		return writeApiError(c, http.StatusInternalServerError, "internal server error")
	}

	hashedValue := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
	state, err := gonanoid.Generate("abcdefghijklmnopqrstwxyzABCDEFGHIJKLMNOPQRSTWXYZ", 32)
	if err != nil {
		logger.Error("Failed to generate state id", zap.Error(err))
		return writeApiError(c, http.StatusInternalServerError, "internal server error")
	}

	// get expiration time
	expiresAt := time.Now().Add(time.Second * 30)

	_, err = database.DB().Exec("INSERT INTO challenges (user_id, state, value, expires_at) VALUES ($1, $2, $3,  $4);", claims.UserId, state, hashedValue, expiresAt)
	if err != nil {
		logger.Error("Failed to insert challenge into database", zap.Error(err))
		return writeApiError(c, http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, GetChallangeResp{
		State:     state,
		Challange: hashedValue,
		ExpiresAt: expiresAt,
	})
}
