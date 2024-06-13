// This module contains all the code for bento routes
package router

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// SetupBentoRoutes is an init function that adds bento routes in a router group.
func SetupBentoRoutes(e RouteGroup) {
	e.POST("/bento/prep", handlePrepBento, useJwtAuth(JWT_ACCESS_TOKEN_TYPE), useValidateRequestBody(prepBentoRequest{}))
	e.PUT("/bento/add-ingridient", handleAddIngridient, useValidateRequestBody(addIngridientRequest{}))
	e.GET("/bento/challenge", handleGetChallenge)
	// Getting a bento only need query parameters
	e.GET("/bento", handleGetBento)
}

func handlePrepBento(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	claims, ok := c.Get("claims").(*jwtAuthClaims)
	if !ok {
		logger.Error("Invalid type casting when getting auth claims.", zap.Any("claims", c.Get("claims")), zap.String(echo.HeaderXRequestID, requestId))
		return writeApiErrorJSON(c, requestId)
	}

	// request body
	body, ok := c.Get("body").(*prepBentoRequest)
	if !ok {
		logger.Error("Invalid type casting when getting request body.", zap.Any("body", c.Get("body")), zap.String(echo.HeaderXRequestID, requestId))
		return writeApiErrorJSON(c, requestId)
	}

	bentoId, err := store.PrepBento(body.Name, claims.UserId, body.PubKey)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code.Name() == store.PG_ERR_UNIQUE_VIOLATION {
			return c.JSON(
				http.StatusBadRequest,
				apiResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "There is another bento you have prepared before with the same name.",
				},
			)
		}
		logger.Error("Failed to prep new bento.", zap.Error(err), zap.String(echo.HeaderXRequestID, requestId))
		return writeApiErrorJSON(c, requestId)
	}

	return c.JSON(
		http.StatusCreated,
		prepBentoResponse{BentoId: bentoId},
	)
}

func handleGetChallenge(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	bentoId := c.QueryParam("bento_id")
	if bentoId == "" {
		return c.JSON(
			http.StatusBadRequest,
			apiResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Missing required query parameter: 'bento_id'",
			},
		)
	}

	challengeId, challenge, err := store.NewChallenge(bentoId)
	if err != nil {
		logger.Error("Failed to get new challenge", zap.Error(err), zap.String("request_id", requestId))
		return writeApiErrorJSON(c, requestId)
	}

	return c.JSON(
		http.StatusCreated,
		getChallengeResponse{
			ChallengeId: challengeId,
			Challenge:   challenge,
		},
	)
}

func handleAddIngridient(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	sugar := zap.L().Sugar()

	claims, err := useJWT(c, JWT_ACCESS_TOKEN_TYPE)
	if err != nil {
		zap.L().Error("Failed to get jwt.", zap.Error(err), zap.String("request_id", requestId))
		return c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	sugar.Info(claims)

	body, ok := c.Get("body").(addIngridientRequest)
	if !ok {
		zap.L().Error("Failed to cast body.", zap.String("request_id", requestId))
		return writeApiErrorJSON(c, requestId)
	}

	// gotta get the pub key from the bento first
	bento, err := store.GetBento(body.BentoId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(
				http.StatusNotFound,
				apiResponse{
					StatusCode: http.StatusNotFound,
					Message:    fmt.Sprintf("Bento with id %s not found.", body.BentoId),
				},
			)
		}
		zap.L().Error("Failed to get bento", zap.String("bento_id", body.BentoId), zap.Error(err))
		return writeApiErrorJSON(c, requestId)
	}
	sugar.Info(bento)

	// verify signed challenge

	err = store.AddIngridient(body.BentoId, body.Key, body.Value)
	if err != nil {
		return err
	}

	return nil
}

func handleGetBento(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	challengeId := c.QueryParam("challenge_id")
	signature := c.QueryParam("signature")
	bentoId := c.QueryParam("bento_id")
	if challengeId == "" || signature == "" || bentoId == "" {
		return c.JSON(
			http.StatusBadRequest,
			apiResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Missing one or more required query parameter(s) 'challenge_id', 'signature', 'bento_id'",
			},
		)
	}

	// get the bento
	bento, err := store.GetBento(bentoId)
	if err != nil {
		zap.L().Error("Failed to get bento.", zap.Error(err), zap.String("bento_id", bentoId), zap.String("request_id", requestId))
		return c.JSON(
			http.StatusUnauthorized,
			apiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    http.StatusText(http.StatusUnauthorized),
			},
		)
	}

	// get the hashed challenge
	challenge, err := store.GetChallenge(challengeId)
	if err != nil {
		zap.L().Error("Failed to get challenge.", zap.Error(err), zap.String("challenge_id", challengeId), zap.String("request_id", requestId))
		return c.JSON(
			http.StatusUnauthorized,
			apiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    http.StatusText(http.StatusUnauthorized),
			},
		)
	}

	if time.Now().After(challenge.ExpiresAt) {
		zap.L().Error("Expired challenge hash.", zap.String("challenge_id", challengeId), zap.String("request_id", requestId))
		// use a go routine because the deletion does not matter to the overall request
		go func(id string) {
			err := store.DeleteChallenge(id)
			if err != nil {
				zap.L().Error("Failed to delete expired challenge.", zap.Error(err))
			}
		}(challengeId)
		return c.JSON(
			http.StatusUnauthorized,
			apiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Usage of expired challenge. Get a new challenge and try again.",
			},
		)
	}

	// verify the signature

	return nil
}
