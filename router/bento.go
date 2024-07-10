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
	e.POST("/bento/prepare", handleNewBento, middleware.Protect())
	e.DELETE("/bento/delete/:bentoId", handleDeleteBento, middleware.Protect())
}

// handleNewBento handles incoming requests to create a new bento. This route must be protected so that no anonymous client can access the api.
func handleNewBento(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	body := new(newBentoReqBody)

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Binding new bento request body.")
	err := c.Bind(body)
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to bind new bento request body.")
		return err
	}

	log.Info().Str(echo.HeaderXRequestID, requestId).Msg("Validating new bento request body.")
	err = c.Validate(body)
	if err != nil {
		c.Set(err_msg_logger_key, "Error when validating new bento request body.")
		return err
	}

	claims, ok := c.Get(middleware.JWT_CLAIMS).(*jwt.JwtClaims)
	if !ok {
		c.Set(err_msg_logger_key, "Failed to cast middleware.JWT_CLAIMS.")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	user, err := store.GetUserWithId(claims.UserId)
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to get user.")
		return err
	}

	if !user.EmailVerified {
		c.Set(err_msg_logger_key, "Aborting creating new bento because user's email has not been verified.")
		c.Set(public_err_msg_key, "Please verify your email first.")
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	bento, err := store.NewBento(body.Name, user.Id, body.PubKey)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code.Name() == store.PG_ERR_UNIQUE_VIOLATION {
			c.Set(public_err_msg_key, fmt.Sprintf("Bento with name '%s' already exists.", body.Name))
			c.Set(err_msg_logger_key, "Aborting new bento creation due to duplication.")
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		c.Set(err_msg_logger_key, "Failed to create new bento.")
		return err
	}
	log.Info().Str(echo.HeaderXRequestID, requestId).Str("bento_name", bento.Name).Str("bento_id", bento.Id).Msg("New bento created.")

	return writeJSON(http.StatusCreated, c, map[string]string{
		"message": "New bento created! Start add ingridients to your bento.",
	})
}

// handleDeleteBento handles incoming requests to delete a bento
func handleDeleteBento(c echo.Context) error {
	bentoId := c.Param("bentoId")
	if bentoId == "" {
		c.Set(err_msg_logger_key, "Missing path param bentoId. This should be impossible to match.")
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if !util.IsValidUUIDv4(bentoId) {
		msg := "Invalid UUID"
		c.Set(err_msg_logger_key, msg)
		c.Set(public_err_msg_key, msg)
		return echo.NewHTTPError(http.StatusNotFound)
	}

	requestId := c.Request().Header.Get(echo.HeaderXRequestID)

	claims, ok := c.Get(middleware.JWT_CLAIMS).(*jwt.JwtClaims)
	if !ok {
		c.Set(err_msg_logger_key, "Failed to cast jwt claims.")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	bento, err := store.GetBentoWithId(bentoId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Set(public_err_msg_key, "Bento not found.")
			c.Set(err_msg_logger_key, "Bento not found.")
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		c.Set(err_msg_logger_key, "Failed to get bento to delete.")
		return err
	}

	// verify if the requesting user is the owner of the bento
	if bento.OwnerId != claims.UserId {
		c.Set(err_msg_logger_key, "Requesting user does not own bento. Aborting deletion.")
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	tx, err := store.StartTx()
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to start tx to delete bento.")
		return err
	}
	_, err = bento.Delete(tx)
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to delete bento.")
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		c.Set(err_msg_logger_key, "Failed to commit transaction to delete bento.")
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Str(echo.HeaderXRequestID, requestId).Msg("Failed to rollback.")
		}
		return err
	}

	return writeJSON(
		http.StatusOK,
		c,
		basicRespBody{
			Msg: "Bento deleted.",
		},
	)
}
