package router

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
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
	e.GET("/bento/order/:bentoId", handleOrderBento)
	e.POST("/bento/prepare", handlePrepareBento, middleware.Protect(), middleware.StructType(reflect.TypeOf(newBentoReqBody{})))
	e.DELETE("/bento/throw/:bentoId", handleThrowBento, middleware.Protect())
	e.POST("/bento/add/ingridients", handleAddIngridients, middleware.Protect(), middleware.StructType(reflect.TypeOf(addIngridientsReqBody{})))
	e.PATCH("/bento/ingridient/rename", handleRenameIngridient, middleware.Protect(), middleware.StructType(reflect.TypeOf(renameIngridientReqBody{})))
	e.PATCH("/bento/ingridient/reseason", handleReseasonIngridient, middleware.Protect(), middleware.StructType(reflect.TypeOf(reseasonIngridientReqBody{})))
	e.DELETE("/bento/ingridient", handleDeleteIngridient, middleware.Protect(), middleware.StructType(reflect.TypeOf(deleteIngridientReqBody{})))
	e.PATCH("/bento/rename", handleRenameBento, middleware.Protect(), middleware.StructType(reflect.TypeOf(renameBentoReqBody{})))
	e.POST("/bento/edit/allow", handleAllowEditBento, middleware.Protect(), middleware.StructType(reflect.TypeOf(allowEditBentoReqBody{})))
	e.PATCH("/bento/edit/revoke", handleRevokeEditBento, middleware.Protect(), middleware.StructType(reflect.TypeOf(revokeShareBentoReqBody{})))
}

// handleOrderBento handles incoming requests to get an existing bento.
func handleOrderBento(c echo.Context) error {
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
		c.Response().Header().Add(echo.HeaderXRequestID, requestId)
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

	entries, err := store.GetEntriesForBento(bento.Id)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get entries for bento.",
			RequestId: requestId,
		}
	}

	ingridients := make([]Ingridient, len(entries))
	for i, e := range entries {
		ingridients[i] = Ingridient{
			Name:  e.Name,
			Value: e.Value,
		}
	}

	resBody := map[string]any{
		"message":     "Here is your bento order.",
		"ingridients": ingridients,
		"request_id":  requestId,
	}

	return writeJSON(http.StatusOK, c, resBody)
}

// handlePrepareBento handles incoming requests to create a new bento. This route must be protected so that no anonymous client can access the api.
func handlePrepareBento(c echo.Context) error {
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

	// need to start a transaction to allow creating perms and bento at the same time
	tx, err := store.StartTx()
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to start transaction to preprare new bento.",
			RequestId: requestId,
		}
	}

	bento, err := store.NewBentoTx(tx, body.Name, user.Id, body.PubKey)
	if err != nil {
		store.Rollback(tx, requestId)
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

	perms, err := store.NewBentoPermissionTx(tx, user.Id, bento.Id, store.O_WRITE|store.O_SHARE|store.O_GRANT_SHARE|store.O_DELETE)
	if err != nil {
		store.Rollback(tx, requestId)
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to create bento permissions for newly prepared bento.",
			RequestId: requestId,
		}
	}
	log.Info().Str(echo.HeaderXRequestID, requestId).Int64("bento_permission_id", perms.Id).Str("bento_id", bento.Id).Msg("New bento permission created.")

	if err := tx.Commit(); err != nil {
		store.Rollback(tx, requestId)
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to commit changes when preparing a new bento.",
			RequestId: requestId,
		}
	}

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
			"message":    "New bento created and ingridients added.",
			"bento_id":   bento.Id,
			"request_id": requestId,
		})
	}

	return writeJSON(http.StatusCreated, c, map[string]string{
		"message":    "New bento created! Start add ingridients to your bento.",
		"bento_id":   bento.Id,
		"request_id": requestId,
	})
}

// handleThrowBento handles incoming requests to delete a bento from the database.
func handleThrowBento(c echo.Context) error {
	bentoId := c.Param("bentoId")
	if bentoId == "" || !util.IsValidUUIDv4(bentoId) {
		return writeJSON(http.StatusNotFound, c, basicRespBody{
			Msg: http.StatusText(http.StatusNotFound),
		})
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
			return apiError{
				Code:      http.StatusNotFound,
				Msg:       "No bento found to throw.",
				PublicMsg: "No bento found to throw.",
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to get bento to delete.",
			Err:       err,
			RequestId: requestId,
		}
	}

	// get perms
	perms, err := store.GetBentoPermissionByUserBentoId(claims.UserId, bento.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusUnauthorized,
				Msg:       "No bento permissions found.",
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get bento permissions.",
			RequestId: requestId,
		}
	}

	// verify if the requesting user is the owner of the bento
	if perms.Permissions&store.O_DELETE == 0 {
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
			Msg:       "Bento deleted.",
			RequestId: requestId,
		},
	)
}

// handleAddIngridients handles incoming requests to add a new ingridient (entry) to a bento
func handleAddIngridients(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	body := new(addIngridientsReqBody)
	if err := readRequestBody(c, body); err != nil {
		return apiError{
			Code:      http.StatusBadRequest,
			Err:       err,
			Msg:       "Failed to read request body.",
			PublicMsg: "Invalid body.",
			RequestId: requestId,
		}
	}
	bento, err := store.GetBentoWithId(body.BentoId)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusNotFound,
				Err:       err,
				Msg:       "Bento not found.",
				PublicMsg: "Bento not found.",
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			RequestId: requestId,
		}
	}

	// verify challenge and signature
	if err := bento.VerifySignature(body.Signature, body.Challenge); err != nil {
		return apiError{
			Code:      http.StatusForbidden,
			Msg:       "Invalid signature.",
			PublicMsg: "Invalid signature.",
			Err:       err,
			RequestId: requestId,
		}
	}

	// jwt claims
	claims, err := middleware.GetJwtClaimsFromContext(c)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get jwt claims from context.",
			RequestId: requestId,
		}
	}

	// get user permissions
	perms, err := store.GetBentoPermissionByUserBentoId(claims.UserId, bento.Id)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get bento permissions.",
			RequestId: requestId,
		}
	}

	// confirm perms to add ingridient
	if perms.Permissions&(store.O_WRITE|store.O_WRITE_INGRIDIENT) == 0 {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       fmt.Sprintf("Requesting user (%s) does not have permission to write ingridient to bento with id '%s'", claims.UserId, bento.Id),
			RequestId: requestId,
		}
	}

	entries := make([]store.BentoEntry, len(body.Ingridients))
	for i := 0; i < len(body.Ingridients); i++ {
		entries[i] = store.NewBentoEntry(body.Ingridients[i].Name, body.Ingridients[i].Value, bento.Id)
	}
	err = store.SaveBentoEntryBatch(entries)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to save batch of entries.",
			RequestId: requestId,
		}
	}
	return writeJSON(http.StatusOK, c, map[string]string{"message": "Ingridients added."})
}

// Handle incoming requests to rename a bento.
func handleRenameBento(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)

	body := new(renameBentoReqBody)
	if err := c.Bind(body); err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to bind rename bento request body.",
			RequestId: requestId,
			Err:       err,
		}
	}
	if err := c.Validate(body); err != nil {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Error when validating rename bento request body.",
			PublicMsg: "Invalid request body",
			RequestId: requestId,
			Err:       err,
		}
	}

	claims, err := middleware.GetJwtClaimsFromContext(c)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get claims from context",
			RequestId: requestId,
		}
	}

	user, err := store.GetUserWithId(claims.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			// no user, so we just return unauthorized
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get user",
			RequestId: requestId,
		}
	}

	bento, err := store.GetBentoWithId(body.BentoId)
	if err != nil {
		if err == sql.ErrNoRows {
			return apiError{
				Code:      http.StatusNotFound,
				Err:       err,
				Msg:       "No bento found with given id.",
				RequestId: requestId,
			}
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get bento with given id.",
			RequestId: requestId,
		}
	}

	perms, err := store.GetBentoPermissionByUserBentoId(user.Id, bento.Id)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get bento permissions.",
			RequestId: requestId,
		}
	}

	if perms.Permissions&(store.O_WRITE|store.O_RENAME_BENTO) == 0 {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       "Requesting user is not owner of bento. Aborting update.",
			RequestId: requestId,
		}
	}

	if err := bento.Rename(body.NewName); err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to rename bento",
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, basicRespBody{
		Msg:       fmt.Sprintf("Bento with id: '%s' renamed to '%s'.", bento.Id, bento.Name),
		RequestId: requestId,
	})
}

// handleRenameIngridient handles incoming requests to rename ingridient(s).
func handleRenameIngridient(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	path := c.Request().URL.Path
	body := new(renameIngridientReqBody)
	if err := readRequestBody(c, body); err != nil {
		return apiError{
			Msg:       "Failed to read request body",
			Path:      path,
			Err:       err,
			RequestId: requestId,
		}
	}
	claims, err := middleware.GetJwtClaimsFromContext(c)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Msg:       "Failed to get jwt claims from context",
			Path:      path,
			Err:       err,
			RequestId: requestId,
		}
	}
	bento, err := store.GetBentoWithId(body.BentoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid bento id. Bento does not exists.")
		}
		return apiError{
			Msg:       "Failed to get bento",
			Path:      path,
			Err:       err,
			RequestId: requestId,
		}
	}
	if err := bento.VerifySignature(body.Challenger.Signature, body.Challenger.Challenge); err != nil {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       "Challenger failed",
			PublicMsg: "Challenger failed",
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}

	bentoPerms, err := store.GetBentoPermissionByUserBentoId(claims.UserId, bento.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apiError{
				Code:      http.StatusUnauthorized,
				Msg:       "No bento permissions for requesting user found",
				PublicMsg: "Invalid credentials. Make sure to have a valid access token.",
				Path:      path,
				Err:       err,
				RequestId: requestId,
			}
		}
		return apiError{
			Msg:       "Failed to get bento permissions for requesting user",
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}

	if bentoPerms.Permissions&(store.O_WRITE|store.O_RENAME_INGRIDIENT|store.O_WRITE_INGRIDIENT) == 0 {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       fmt.Sprintf("Requesting user does not have permissions to rename ingridient. Bento ID: %s", bento.Id),
			PublicMsg: fmt.Sprintf("You do not have permissions to rename ingridients in bento with ID: %s", bento.Id),
			Path:      path,
			Err:       err,
			RequestId: requestId,
		}
	}

	if err := store.RenameIngridient(bento.Id, body.OldName, body.NewName); err != nil {
		return apiError{
			Msg:       "Failed to rename ingridient",
			Path:      path,
			Err:       err,
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, basicRespBody{Msg: "Ingridients renamed", RequestId: requestId})
}

func handleReseasonIngridient(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	path := c.Request().URL.Path
	body := new(reseasonIngridientReqBody)
	if err := readRequestBody(c, body); err != nil {
		code := http.StatusInternalServerError
		if _, ok := err.(validator.ValidationErrors); ok {
			code = http.StatusBadRequest
		}
		return apiError{
			Code:      code,
			Msg:       "Failed to read request body",
			Path:      path,
			Err:       err,
			RequestId: requestId,
		}
	}
	claims, err := middleware.GetJwtClaimsFromContext(c)
	if err != nil {
		return apiError{
			Msg:       "Failed to get jwt claims from context",
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}
	bento, err := store.GetBentoWithId(body.BentoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid bento id. Bento does not exists.")
		}
		return apiError{
			Msg:       "Failed to get bento",
			Path:      path,
			Err:       err,
			RequestId: requestId,
		}
	}
	if err := bento.VerifySignature(body.Challenger.Signature, body.Challenger.Challenge); err != nil {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       "Challenger failed",
			PublicMsg: "Challenger failed",
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}
	bentoPerms, err := store.GetBentoPermissionByUserBentoId(claims.UserId, bento.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apiError{
				Code:      http.StatusUnauthorized,
				Msg:       "No bento permissions for requesting user found",
				PublicMsg: "Invalid credentials. Make sure to have a valid access token.",
				Path:      path,
				Err:       err,
				RequestId: requestId,
			}
		}
		return apiError{
			Msg:       "Failed to get bento permissions for requesting user",
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}
	if bentoPerms.Permissions&(store.O_WRITE|store.O_WRITE_INGRIDIENT) == 0 {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       fmt.Sprintf("Requesting user does not have permissions to re-season ingridient. Bento ID: %s", bento.Id),
			PublicMsg: fmt.Sprintf("You do not have permissions to re-season ingridients in bento with ID: %s", bento.Id),
			Path:      path,
			Err:       err,
			RequestId: requestId,
		}
	}

	if err := store.ReseasonIngridient(bento.Id, body.Name, body.Value); err != nil {
		return apiError{
			Msg:       "Failed to change ingridient value",
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, basicRespBody{Msg: "Bento ingridient re-seasoned.", RequestId: requestId})
}

func handleDeleteIngridient(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	path := c.Request().URL.Path
	body := new(deleteIngridientReqBody)
	if err := readRequestBody(c, body); err != nil {
		code := http.StatusInternalServerError
		if _, ok := err.(validator.ValidationErrors); ok {
			code = http.StatusBadRequest
		}
		return apiError{
			Code:      code,
			Err:       err,
			Path:      path,
			RequestId: requestId,
			Msg:       "Failed to read request body.",
		}
	}
	claims, err := middleware.GetJwtClaimsFromContext(c)
	if err != nil {
		return apiError{
			Err:       err,
			Path:      path,
			RequestId: requestId,
			Msg:       "Failed to get jwt claims from context",
		}
	}
	bento, err := store.GetBentoWithId(body.BentoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apiError{
				Code:      http.StatusBadRequest,
				Err:       err,
				Path:      path,
				RequestId: requestId,
				PublicMsg: fmt.Sprintf("No bento with ID: %s found", body.BentoId),
			}
		}
		return apiError{
			Err:       err,
			Path:      path,
			RequestId: requestId,
			Msg:       "Failed to get bento",
		}
	}
	if err := bento.VerifySignature(body.Challenger.Signature, body.Challenger.Challenge); err != nil {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       "Challenger failed",
			PublicMsg: "Challenger failed",
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}
	bentoPerms, err := store.GetBentoPermissionByUserBentoId(claims.UserId, bento.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apiError{
				Code:      http.StatusUnauthorized,
				Msg:       "No bento permissions for requesting user found",
				PublicMsg: "Invalid credentials. Make sure to have a valid access token.",
				Path:      path,
				Err:       err,
				RequestId: requestId,
			}
		}
		return apiError{
			Msg:       "Failed to get bento permissions for requesting user",
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}
	if bentoPerms.Permissions&(store.O_DELETE|store.O_DELETE_INGRIDIENT) == 0 {
		return apiError{
			Code:      http.StatusUnauthorized,
			Msg:       "Requesting user has no permissions to delete ingridient",
			PublicMsg: fmt.Sprintf("You do not have permissions to delete ingridients in bento with ID: %s", bento.Id),
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}
	if err := store.DeleteIngridient(bento.Id, body.Name); err != nil {
		return apiError{
			Msg:       "Failed to delete ingridient from bento",
			Err:       err,
			Path:      path,
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, basicRespBody{Msg: "Bento ingridient removed.", RequestId: requestId})
}

func handleAllowEditBento(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	body := new(allowEditBentoReqBody)

	if err := c.Bind(body); err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to bind body",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	if err := c.Validate(body); err != nil {
		return apiError{
			Code:      http.StatusBadRequest,
			Msg:       "Error when validating body",
			PublicMsg: "Invalid request body",
			Err:       err,
			RequestId: requestId,
			Path:      c.Request().URL.Path,
		}
	}

	claims, err := middleware.GetJwtClaimsFromContext(c)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get claims from context.",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	// varify if target user exists or not
	targetUser, err := store.GetUserWithEmail(body.ShareToEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "Target email does not belong to any user.")
		}
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to get target user.",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	bento, err := store.GetBentoWithId(body.BentoId)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No bento found with id: %s", body.BentoId))
		}
		return apiError{
			Code:      http.StatusInsufficientStorage,
			Err:       err,
			Msg:       "Failed to get bento.",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	exists, err := store.ExistsBentoPermissionByUserBentoId(targetUser.Id, bento.Id)
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to check if bento permission exists.",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	} else if exists {
		return writeJSON(http.StatusOK, c, basicRespBody{
			Msg:       "Bento has been previously shared to user. Refer to 'https://github.com/juancwu/konbini/blob/main/.github/docs/DOCUMENTATION.md#share-bento' for more information.",
			RequestId: requestId,
		})
	}

	if err := bento.VerifySignature(body.Signature, body.Challenge); err != nil {
		return apiError{
			Code:      http.StatusUnauthorized,
			Err:       err,
			Msg:       "Failed to verify signature.",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	perms, err := store.GetBentoPermissionByUserBentoId(claims.UserId, bento.Id)
	if err != nil {
		return apiError{
			Code:      http.StatusUnauthorized,
			Err:       err,
			Msg:       "Failed to get bento permissions.",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	if perms.Permissions&store.O_SHARE == 0 {
		return apiError{
			Code:      http.StatusUnauthorized,
			Err:       err,
			Msg:       "No permission to share bento",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	tx, err := store.StartTx()
	if err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to start transaction.",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	// It can only grant up to their own level of permission
	targetUserPerms := store.O_NO_PERMS
	if body.PermissionLevels != nil && len(body.PermissionLevels) > 0 {
		for _, level := range body.PermissionLevels {
			if level == store.S_ALL {
				// remove the grant share bit from the permissions if they had any
				// by default, granting the ability to share should be explicitly
				// given in another route
				targetUserPerms = perms.Permissions & (^store.O_GRANT_SHARE)
				// no need to continue since we got all the perms we need
				break
			}
			// assigned the bit to the target user permissions
			oLevel, ok := store.TextToBinPerms[level]
			if !ok {
				continue
			}
			targetUserPerms = targetUserPerms | (perms.Permissions & oLevel)
		}
	}

	if _, err := store.NewBentoPermissionTx(tx, targetUser.Id, bento.Id, targetUserPerms); err != nil {
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to create bento permissions.",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	if err := tx.Commit(); err != nil {
		store.Rollback(tx, requestId)
		return apiError{
			Code:      http.StatusInternalServerError,
			Err:       err,
			Msg:       "Failed to commit bento permission changes.",
			Path:      c.Request().URL.Path,
			RequestId: requestId,
		}
	}

	return writeJSON(http.StatusOK, c, basicRespBody{Msg: fmt.Sprintf("Bento shared with '%s'", targetUser.Email), RequestId: requestId})
}

func handleRevokeEditBento(c echo.Context) error {
	requestId := c.Request().Header.Get(echo.HeaderXRequestID)
	// path := c.Request().URL.Path

	return writeJSON(http.StatusOK, c, basicRespBody{
		Msg:       "Share revoked.",
		RequestId: requestId,
	})
}
