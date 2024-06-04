package router

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/middleware"
	bentomodel "github.com/juancwu/konbini/server/models/bento"
	entrymodel "github.com/juancwu/konbini/server/models/entry"
	usermodel "github.com/juancwu/konbini/server/models/user"
	"github.com/juancwu/konbini/server/service"
	"github.com/juancwu/konbini/server/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func validateUUID(uuid string) error {
	if utils.IsValidUUIDV4(uuid) {
		return nil
	}
	return fmt.Errorf("The given id is not a proper UUID v4: %s", uuid)
}

func SetupBentoRoutes(e RouteGroup) {
	e.GET("/bento/menu", handleListPersonalBentos, middleware.JwtAuthMiddleware)
	e.DELETE("/bento/toss/:id", handleDeletePersonalBento, middleware.JwtAuthMiddleware, ValidateRequest(
		ValidatorOptions{
			Field:    "id",
			From:     VALIDATE_PARAM,
			Required: true,
			Validate: validateUUID,
		},
	))
	e.PATCH("/bento/rebrand/:id", handleRebrandBento, ValidateRequest(
		ValidatorOptions{
			Field:    "id",
			From:     VALIDATE_PARAM,
			Required: true,
			Validate: validateUUID,
		},
	))

	// New api routes
	e.POST("/bento/prep", handleNewPersonalBento, middleware.JwtAuthMiddleware)
	e.GET("/bento/:id", handleGetBento, ValidateRequest(
		ValidatorOptions{
			Field:    "id",
			From:     VALIDATE_PARAM,
			Required: true,
			Validate: validateUUID,
		},
		ValidatorOptions{
			Field:    "X-Bento-Hashed",
			From:     VALIDATE_HEADER,
			Required: true,
		},
		ValidatorOptions{
			Field:    "X-Bento-Signature",
			From:     VALIDATE_HEADER,
			Required: true,
		},
	))
}

func handleGetBento(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	bentoId := c.Param("id")
	// need to query the bento first to get the public key
	bento, err := bentomodel.GetPersonalBento(bentoId)
	if err != nil {
		logger.Error("Failed to get personal bento", zap.Error(err))
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Personal bento not found.")
		}

		return c.String(http.StatusInternalServerError, "Failed to get personal bento.")
	}
	hashed := c.Request().Header.Get("X-Bento-Hashed")
	signature := c.Request().Header.Get("X-Bento-Signature")

	decodedHashed, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		logger.Error("Failed to decode base64 hashed challenge", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to decode hashed challenge")
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		logger.Error("Failed to decode base64 signature", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to decode signature")
	}

	err = service.VerifyBentoSignature(decodedHashed, decodedSignature, []byte(bento.PubKey))
	if err != nil {
		logger.Error("Failed to verify bento signature", zap.Error(err))
		return c.String(http.StatusUnauthorized, "Invalid signature")
	}

	// get the key-value pairs for the bento
	bento.KeyVals, err = entrymodel.GetPersonalBentoEntries(bento.Id)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get personal bento entries", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to get personal bento.")
	}

	return c.JSON(http.StatusOK, bento)
}

type NewPersonalBentoReqBody struct {
	Name       string   `json:"name" validate:"required,min=1"`
	PublickKey string   `json:"public_key" validate:"required,min=1"`
	KeyVals    []string `json:"keyvals" validate:"required,ValidateStringSlice"`
}

func handleNewPersonalBento(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	reqBody := new(NewPersonalBentoReqBody)

	if err := c.Bind(reqBody); err != nil {
		logger.Error("Failed to bind request body", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Personal bento service down. Please try again later.")
	}
	if err := c.Validate(reqBody); err != nil {
		logger.Error("Request body validation failed", zap.Error(err))
		return c.String(http.StatusBadRequest, "Bad request")
	}

	claims := c.Get("claims").(*service.JwtCustomClaims)

	// get user
	isReal, err := usermodel.IsRealUser(claims.UserId)
	if err != nil {
		logger.Error("Failed to get user", zap.Error(err))
		return c.String(http.StatusBadRequest, "You must be an existing member of Konbini to create personal bentos.")
	}
	if !isReal {
		logger.Error("User with id in claims returned as not real. Possible old active access token used.")
		return c.String(http.StatusBadRequest, "You must be an existing member of Konbini to create personal bentos.")
	}

	// check if user has the same personal bento
	exists, err := bentomodel.PersonalBentoExistsWithName(claims.UserId, reqBody.Name)
	if err != nil {
		logger.Error("Failed to check if user has personal bento with same name", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Personal bento service down. Please try again later.")
	}

	if exists {
		logger.Error("Attempt to create a new personal bento with the same name.")
		return c.String(http.StatusBadRequest, "Another personal bento with the same name already exists. If you wish to replace the bento, please delete it and create a new one.")
	}

	// create a new transaction
	tx, err := database.DB().Begin()
	if err != nil {
		logger.Error("Failed to begin transaction", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Personal bento servide down. Please try again later.")
	}

	bentoId, err := bentomodel.NewPersonalBento(tx, claims.UserId, reqBody.Name, reqBody.PublickKey)
	if err != nil {
		logger.Error("Failed to create new personal bento", zap.Error(err))
		if err := tx.Rollback(); err != nil {
			logger.Error("Failed to rollback transaction", zap.Error(err))
		}
		return c.String(http.StatusInternalServerError, "Failed to create personal bento. Please try again later.")
	}
	logger.Info("New personal bento created.", zap.String("user_id", claims.UserId), zap.String("bento_id", bentoId))
	logger.Info("Registering personal bento entries...")

	err = entrymodel.CreateEntries(tx, bentoId, reqBody.KeyVals)
	if err != nil {
		logger.Error("Failed to register personal bento entries", zap.Error(err))
		if err := tx.Rollback(); err != nil {
			logger.Error("Failed to rollback transaction", zap.Error(err))
		}
		return c.String(http.StatusInternalServerError, "Failed to create personal bento. Please try again later.")
	}
	logger.Info("Personal bento entries registered.", zap.String("bento_id", bentoId))

	err = tx.Commit()
	if err != nil {
		logger.Error("Failed to commit transaction", zap.Error(err))
		if err := tx.Rollback(); err != nil {
			logger.Error("Failed to rollback transaction", zap.Error(err))
		}
		return c.String(http.StatusInternalServerError, "Failed to create personal bento. Please try again later.")
	}

	return c.String(http.StatusCreated, bentoId)
}

func handleListPersonalBentos(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	claims, ok := c.Get("claims").(*service.JwtCustomClaims)
	if !ok {
		logger.Error("No claims found. Needed to get list of personal bentos.")
		return c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	logger.Info("Getting list of personal bentos...")
	bentos, err := bentomodel.ListPersonalBentos(claims.UserId)
	if err != nil {
		logger.Error("Failed to get list of personal bentos", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to get list of personal bentos from database.")
	}

	if bentos == nil {
		return c.JSON(http.StatusOK, []bentomodel.PersonalBento{})
	}

	return c.JSON(http.StatusOK, bentos)
}

func handleDeletePersonalBento(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	claims, ok := c.Get("claims").(*service.JwtCustomClaims)
	if !ok {
		logger.Error("No claims found. Needed to delete bento.")
		return c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	id := c.Param("id")
	bento, err := bentomodel.GetPersonalBento(id)
	if err != nil {
		logger.Error("Failed to get personal bento", zap.Error(err))
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Personal bento not found.")
		}
		return c.String(http.StatusInternalServerError, "Failed to get personal bento.")
	}

	hashed := c.Request().Header.Get("X-Bento-Hashed")
	signature := c.Request().Header.Get("X-Bento-Signature")

	decodedHashed, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		logger.Error("Failed to decode base64 hashed challenge", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to decode hashed challenge")
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		logger.Error("Failed to decode base64 signature", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to decode signature")
	}

	err = service.VerifyBentoSignature(decodedHashed, decodedSignature, []byte(bento.PubKey))
	if err != nil {
		logger.Error("Failed to verify bento signature", zap.Error(err))
		return c.String(http.StatusUnauthorized, "Invalid signature")
	}

	deleted, err := bentomodel.DeletePersonalBento(claims.UserId, bento.Id)
	if err != nil || !deleted {
		logger.Error("Failed to delete personal bento", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to delete perosnal bento.")
	}

	return c.String(http.StatusOK, "Deleted")
}

const (
	UPDATE_ACTION_ADD    = "add"
	UPDATE_ACTION_DELETE = "delete"
	UPDATE_ACTION_UPDATE = "update"
)

func handleRebrandBento(c echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	id := c.Param("id")

	reqBody := new(RebrandBentoRequestBody)
	if err := c.Bind(reqBody); err != nil {
		logger.Error("Failed to bind request body", zap.Error(err))
		return writeApiError(c, http.StatusBadRequest, "bad request")
	}
	if err := c.Validate(reqBody); err != nil {
		logger.Error("Requests body validation failed", zap.Error(err))
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return writeApiReqBodyError(c, http.StatusBadRequest, ve)
		}
		return writeApiError(c, http.StatusBadRequest, "bad request")
	}

	// TODO: add check to verify if user has permission to perform update actions

	// make sure that the bento actually exists
	bento, err := bentomodel.GetPersonalBento(id)
	if err != nil {
		logger.Error("Failed to get personal bento", zap.Error(err))
		if err == sql.ErrNoRows {
			return writeApiError(c, http.StatusNotFound, "Personal bento not found.")
		}
		return writeApiError(c, http.StatusInternalServerError, "Failed to get personal bento.")
	}

	// we still need to authenticate this
	hashed := c.Request().Header.Get("X-Bento-Hashed")
	signature := c.Request().Header.Get("X-Bento-Signature")

	decodedHashed, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		logger.Error("Failed to decode base64 hashed challenge", zap.Error(err))
		return writeApiError(c, http.StatusInternalServerError, "Failed to decode hashed challenge")
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		logger.Error("Failed to decode base64 signature", zap.Error(err))
		return writeApiError(c, http.StatusInternalServerError, "Failed to decode signature")
	}

	err = service.VerifyBentoSignature(decodedHashed, decodedSignature, []byte(bento.PubKey))
	if err != nil {
		logger.Error("Failed to verify bento signature", zap.Error(err))
		return writeApiError(c, http.StatusUnauthorized, "Invalid signature")
	}

	// update the bento's name
	_, err = database.DB().Exec("UPDATE personal_bentos SET name = $1 WHERE id = $2;", reqBody.NewName, id)
	if err != nil {
		logger.Error("Failed to update bento's name", zap.Error(err))
		return writeApiError(c, http.StatusInternalServerError, "internal server error")
	}

	return writeNoBody(c, http.StatusOK)
}
