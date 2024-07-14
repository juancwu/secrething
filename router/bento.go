package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/juancwu/konbini/jwt"
	"github.com/juancwu/konbini/middleware"
	"github.com/juancwu/konbini/store"
	"github.com/juancwu/konbini/util"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// SetupBentoRoutes setups the routes for bento services.
func SetupBentoRoutes(e RouterGroup) {
	e.GET("/bento/:bentoId", handleGetBento)
	e.POST("/bento/prepare", handleNewBento, middleware.Protect())
	e.DELETE("/bento/delete/:bentoId", handleDeleteBento, middleware.Protect())
}

// handleGetBento handles incoming requests to get an existing bento.
func handleGetBento(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	signature := c.QueryParam("signature")
	challenge := c.QueryParam("challenge")
	bentoId := c.Param("bentoId")

	if bentoId == "" || !util.IsValidUUIDv4(bentoId) {
		return c.NoContent(http.StatusNotFound)
	}

	if signature == "" {
		return apiError{
			Code:      http.StatusBadRequest,
			RequestId: requestId,
			Msg:       "Missing requried query parameter 'signature'",
			PublicMsg: "Missing requried query parameter 'signature'. It should be hex encoded.",
		}
	}

	if challenge == "" {
		return apiError{
			Code:      http.StatusBadRequest,
			RequestId: requestId,
			Msg:       "Missing requried query parameter 'challenge'",
			PublicMsg: "Missing requried query parameter 'challenge'. It should be hex encoded.",
		}
	}

	// get bento
	bento, err := store.GetBentoWithId(bentoId)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	// verify signature
	err = bento.VerifySignature(signature, challenge)
	if err != nil {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Invalid signature.",
			PublicMsg: "Invalid signature.",
			Err:       err,
			RequestId: requestId,
		}
	}

	return nil
}

// handleNewBento handles incoming requests to create a new bento. This route must be protected so that no anonymous client can access the api.
func handleNewBento(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	body := new(newBentoReqBody)

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Binding new bento request body.")
	err := c.Bind(body)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to bind new bento request body.",
			RequestId: requestId,
			Err:       err,
		}
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Validating new bento request body.")
	err = c.Validate(body)
	if err != nil {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Error when validating new bento request body.",
			PublicMsg: "Invalid request body",
			RequestId: requestId,
			Err:       err,
		}
	}

	claims, ok := c.Get(middleware.JWT_CLAIMS).(*jwt.JwtClaims)
	if !ok {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to cast middleware.JWT_CLAIMS.",
			RequestId: requestId,
		}
	}

	user, err := store.GetUserWithId(claims.UserId)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to get user.",
			RequestId: requestId,
		}
	}

	if !user.EmailVerified {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       "Aborting creating new bento because user's email has not been verified.",
			PublicMsg: "Please verify your email first.",
			RequestId: requestId,
		}
	}

	bento, err := store.NewBento(body.Name, user.Id, body.PubKey)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code.Name() == store.PG_ERR_UNIQUE_VIOLATION {
			return apiError{
				Code:      http.StatusBadRequest,
				Msg:       "Aborting new bento creation due to duplication.",
				PublicMsg: fmt.Sprintf("Bento with name '%s' already exists.", body.Name),
				Err:       pgErr,
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to create new bento.",
			Err:       err,
			RequestId: requestId,
		}
	}
	log.Info().Str(echo.HeaderXRequestID, requestId).Str("bento_name", bento.Name).Str("bento_id", bento.Id).Msg("New bento created.")

	if body.Ingridients != nil && len(body.Ingridients) > 0 {
		log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Trying to add ingridients.")
		entries := make([]store.BentoEntry, len(body.Ingridients))
		for i, ingridient := range body.Ingridients {
			entries[i] = store.NewBentoEntry(ingridient.Name, ingridient.Value, bento.Id)
		}
		if err := store.SaveBentoEntryBatch(entries); err != nil {
			return apiError{
				Code:      http.StatusOK,
				Err:       err,
				Msg:       "Failed to add ingridients in the same request to prepare bento.",
				PublicMsg: "New bento created, but ingridients were not able to be added.",
				RequestId: requestId,
			}
		}
		return writeJSON(http.StatusCreated, c, map[string]string{
			"message": "New bento created and ingridients added.",
		})
	}

	return writeJSON(http.StatusCreated, c, map[string]string{
		"message": "New bento created! Start add ingridients to your bento.",
	})
}

// handleDeleteBento handles incoming requests to delete a bento
func handleDeleteBento(c echo.Context) error {
	bentoId := c.Param("bentoId")
	if bentoId == "" || !util.IsValidUUIDv4(bentoId) {
		return c.NoContent(http.StatusNotFound)
	}

	requestId := c.Request().Header.Get(echo.HeaderXRequestID)

	claims, ok := c.Get(middleware.JWT_CLAIMS).(*jwt.JwtClaims)
	if !ok {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to cast jwt claims.",
			RequestId: requestId,
		}
	}

	bento, err := store.GetBentoWithId(bentoId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to get bento to delete.",
			Err:       err,
			RequestId: requestId,
		}
	}

	// verify if the requesting user is the owner of the bento
	if bento.OwnerId != claims.UserId {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       "Requesting user does not own bento. Aborting deletion.",
			PublicMsg: http.StatusText(http.StatusUnauthorized),
			RequestId: requestId,
		}
	}

	tx, err := store.StartTx()
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to start tx to delete bento.",
			RequestId: requestId,
			Err:       err,
		}
	}
	_, err = bento.Delete(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to delete bento.",
			RequestId: requestId,
			Err:       err,
		}
	}

	err = tx.Commit()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to commit tx to delete bento.",
			Err:       err,
			RequestId: requestId,
		}
	}

	return writeJSON(
		http.StatusOK,
		c,
		basicRespBody{
			Msg: "Bento deleted.",
		},
	)
}
