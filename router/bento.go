// This module contains all the code for bento routes
package router

import (
	"net/http"

	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// SetupBentoRoutes is an init function that adds bento routes in a router group.
func SetupBentoRoutes(e RouteGroup) {
	e.POST("/bento/prep", handlePrepBento, useJwtAuth(JWT_ACCESS_TOKEN_TYPE), useValidateRequestBody(prepBentoRequest{}))
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
