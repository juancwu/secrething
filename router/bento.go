package router

import (
	"fmt"
	"net/http"

	"github.com/juancwu/konbini/jwt"
	"github.com/juancwu/konbini/middleware"
	"github.com/juancwu/konbini/store"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// SetupBentoRoutes setups the routes for bento services.
func SetupBentoRoutes(e RouterGroup) {
	e.POST("/bento/prepare", handleNewBento, middleware.Protect())
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
