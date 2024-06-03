package router

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/middleware"
	bentomodel "github.com/juancwu/konbini/server/models/bento"
	entrymodel "github.com/juancwu/konbini/server/models/entry"
	usermodel "github.com/juancwu/konbini/server/models/user"
	"github.com/juancwu/konbini/server/service"
	"github.com/juancwu/konbini/server/utils"
	"github.com/labstack/echo/v4"
)

func SetupBentoRoutes(e RouteGroup) {
	// DEPRACATED
	e.POST("/bento/personal/new", handleNewPersonalBento, middleware.JwtAuthMiddleware)
	e.GET("/bento/personal/:id", handleGetPersonalBento)
	e.GET("/bento/personal/list", handleListPersonalBentos, middleware.JwtAuthMiddleware)
	e.DELETE("/bento/personal/:id", handleDeletePersonalBento, middleware.JwtAuthMiddleware)
	e.PATCH("/bento/personal/:id", handleUpdatePersonalBento)

	// New api routes
	e.GET("/bento/:id", handleGetBento, middleware.ValidateRequest(
		middleware.ValidatorOptions{
			Field:    "id",
			From:     middleware.VALIDATE_PARAM,
			Required: true,
			Validate: func(s string) error {
				if utils.IsValidUUIDV4(s) {
					return nil
				}
				return fmt.Errorf("The given id is not a proper UUID v4: %s", s)
			},
		},
		// the two following validators do not need to properly validate the values
		// because we just need to make sure it required so the existence of a non-empty
		// string is already being taken account for
		middleware.ValidatorOptions{
			Field:    "X-Bento-Hashed",
			From:     middleware.VALIDATE_HEADER,
			Required: true,
			Validate: func(s string) error {
				return nil
			},
		},
		middleware.ValidatorOptions{
			Field:    "X-Bento-Signature",
			From:     middleware.VALIDATE_HEADER,
			Required: true,
			Validate: func(s string) error {
				return nil
			},
		},
	))
}

func handleGetBento(c echo.Context) error {
	bentoId := c.Param("id")
	// need to query the bento first to get the public key
	bento, err := bentomodel.GetPersonalBento(bentoId)
	if err != nil {
		utils.Logger().Errorf("Failed to get personal bento: %v\n", err)
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Personal bento not found.")
		}

		return c.String(http.StatusInternalServerError, "Failed to get personal bento.")
	}
	hashed := c.Request().Header.Get("X-Bento-Hashed")
	signature := c.Request().Header.Get("X-Bento-Signature")

	decodedHashed, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		utils.Logger().Errorf("Failed to decode base64 hashed challenge: %s\n", err)
		return c.String(http.StatusInternalServerError, "Failed to decode hashed challenge")
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		utils.Logger().Errorf("Failed to decode base64 signature: %s\n", err)
		return c.String(http.StatusInternalServerError, "Failed to decode signature")
	}

	err = service.VerifyBentoSignature(decodedHashed, decodedSignature, []byte(bento.PubKey))
	if err != nil {
		utils.Logger().Errorf("Failed to verify bento signature: %v\n", err)
		return c.String(http.StatusUnauthorized, "Invalid signature")
	}

	// get the key-value pairs for the bento
	bento.KeyVals, err = entrymodel.GetPersonalBentoEntries(bento.Id)
	if err != nil && err != sql.ErrNoRows {
		utils.Logger().Errorf("Failed to get personal bento entries: %v\n", err)
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
	reqBody := new(NewPersonalBentoReqBody)

	if err := c.Bind(reqBody); err != nil {
		utils.Logger().Errorf("Failed to bind request body: %v\n", err)
		return c.String(http.StatusInternalServerError, "Personal bento service down. Please try again later.")
	}
	if err := c.Validate(reqBody); err != nil {
		utils.Logger().Errorf("Request body validation failed: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad request")
	}

	claims := c.Get("claims").(*service.JwtCustomClaims)

	// get user
	isReal, err := usermodel.IsRealUser(claims.UserId)
	if err != nil {
		utils.Logger().Errorf("Failed to get user: %v\n", err)
		return c.String(http.StatusBadRequest, "You must be an existing member of Konbini to create personal bentos.")
	}
	if !isReal {
		utils.Logger().Error("User with id in claims returned as not real. Possible old active access token used.")
		return c.String(http.StatusBadRequest, "You must be an existing member of Konbini to create personal bentos.")
	}

	// check if user has the same personal bento
	exists, err := bentomodel.PersonalBentoExistsWithName(claims.UserId, reqBody.Name)
	if err != nil {
		utils.Logger().Errorf("Failed to check if user has personal bento with same name: %v\n", err)
		return c.String(http.StatusInternalServerError, "Personal bento service down. Please try again later.")
	}

	if exists {
		utils.Logger().Error("Attempt to create a new personal bento with the same name.")
		return c.String(http.StatusBadRequest, "Another personal bento with the same name already exists. If you wish to replace the bento, please delete it and create a new one.")
	}

	// create a new transaction
	tx, err := database.DB().Begin()
	if err != nil {
		utils.Logger().Errorf("Failed to begin transaction: %v\n", err)
		return c.String(http.StatusInternalServerError, "Personal bento servide down. Please try again later.")
	}

	bentoId, err := bentomodel.NewPersonalBento(tx, claims.UserId, reqBody.Name, reqBody.PublickKey)
	if err != nil {
		utils.Logger().Errorf("Failed to create new personal bento: %v\n", err)
		if err := tx.Rollback(); err != nil {
			utils.Logger().Errorf("Failed to rollback transaction: %v\n", err)
		}
		return c.String(http.StatusInternalServerError, "Failed to create personal bento. Please try again later.")
	}
	utils.Logger().Info("New personal bento created.", "user_id", claims.UserId, "bento_id", bentoId)
	utils.Logger().Info("Registering personal bento entries...")

	err = entrymodel.CreateEntries(tx, bentoId, reqBody.KeyVals)
	if err != nil {
		utils.Logger().Errorf("Failed to register personal bento entries: %v\n", err)
		if err := tx.Rollback(); err != nil {
			utils.Logger().Errorf("Failed to rollback transaction: %v\n", err)
		}
		return c.String(http.StatusInternalServerError, "Failed to create personal bento. Please try again later.")
	}
	utils.Logger().Info("Personal bento entries registered.", "bento_id", bentoId)

	err = tx.Commit()
	if err != nil {
		utils.Logger().Errorf("Failed to commit transaction: %v\n", err)
		if err := tx.Rollback(); err != nil {
			utils.Logger().Errorf("Failed to rollback transaction: %v\n", err)
		}
		return c.String(http.StatusInternalServerError, "Failed to create personal bento. Please try again later.")
	}

	return c.String(http.StatusCreated, bentoId)
}

// TODO: update to use challenge headers middlware
func handleGetPersonalBento(c echo.Context) error {
	id := c.Param("id")
	if !utils.IsValidUUIDV4(id) {
		utils.Logger().Errorf("Invalid uuid when getting personal bento: %s\n", id)
		return c.String(http.StatusBadRequest, "Invalid uuid.")
	}
	bento, err := bentomodel.GetPersonalBento(id)
	if err != nil {
		utils.Logger().Errorf("Failed to get personal bento: %v\n", err)
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Personal bento not found.")
		}

		return c.String(http.StatusInternalServerError, "Failed to get personal bento.")
	}
	keyvals, err := entrymodel.GetPersonalBentoEntries(bento.Id)
	if err != nil && err != sql.ErrNoRows {
		utils.Logger().Errorf("Failed to get personal bento entries: %v\n", err)
		return c.String(http.StatusInternalServerError, "Failed to get personal bento.")
	}
	bento.KeyVals = keyvals

	hashed := c.Request().Header.Get("X-Bento-Hashed")
	signature := c.Request().Header.Get("X-Bento-Signature")

	decodedHashed, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		utils.Logger().Errorf("Failed to decode base64 hashed challenge: %s\n", err)
		return c.String(http.StatusInternalServerError, "Failed to decode hashed challenge")
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		utils.Logger().Errorf("Failed to decode base64 signature: %s\n", err)
		return c.String(http.StatusInternalServerError, "Failed to decode signature")
	}

	err = service.VerifyBentoSignature(decodedHashed, decodedSignature, []byte(bento.PubKey))
	if err != nil {
		utils.Logger().Errorf("Failed to verify bento signature: %v\n", err)
		return c.String(http.StatusUnauthorized, "Invalid signature")
	}

	return c.JSON(http.StatusOK, bento)
}

func handleListPersonalBentos(c echo.Context) error {
	claims, ok := c.Get("claims").(*service.JwtCustomClaims)
	if !ok {
		utils.Logger().Errorf("No claims found. Needed to get list of personal bentos.")
		return c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	utils.Logger().Info("Getting list of personal bentos...")
	bentos, err := bentomodel.ListPersonalBentos(claims.UserId)
	if err != nil {
		utils.Logger().Errorf("Failed to get list of personal bentos: %v\n", err)
		return c.String(http.StatusInternalServerError, "Failed to get list of personal bentos from database.")
	}

	if bentos == nil {
		return c.JSON(http.StatusOK, []bentomodel.PersonalBento{})
	}

	return c.JSON(http.StatusOK, bentos)
}

func handleDeletePersonalBento(c echo.Context) error {
	claims, ok := c.Get("claims").(*service.JwtCustomClaims)
	if !ok {
		utils.Logger().Errorf("No claims found. Needed to delete bento.")
		return c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}

	id := c.Param("id")
	if !utils.IsValidUUIDV4(id) {
		utils.Logger().Errorf("Invalid uuid when deleting personal bento: %s\n", id)
		return c.String(http.StatusBadRequest, "Invalid uuid")
	}

	bento, err := bentomodel.GetPersonalBento(id)
	if err != nil {
		utils.Logger().Errorf("Failed to get personal bento: %v\n", err)
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Personal bento not found.")
		}
		return c.String(http.StatusInternalServerError, "Failed to get personal bento.")
	}

	hashed := c.Request().Header.Get("X-Bento-Hashed")
	signature := c.Request().Header.Get("X-Bento-Signature")

	decodedHashed, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		utils.Logger().Errorf("Failed to decode base64 hashed challenge: %s\n", err)
		return c.String(http.StatusInternalServerError, "Failed to decode hashed challenge")
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		utils.Logger().Errorf("Failed to decode base64 signature: %s\n", err)
		return c.String(http.StatusInternalServerError, "Failed to decode signature")
	}

	err = service.VerifyBentoSignature(decodedHashed, decodedSignature, []byte(bento.PubKey))
	if err != nil {
		utils.Logger().Errorf("Failed to verify bento signature: %v\n", err)
		return c.String(http.StatusUnauthorized, "Invalid signature")
	}

	deleted, err := bentomodel.DeletePersonalBento(claims.UserId, bento.Id)
	if err != nil || !deleted {
		utils.Logger().Errorf("Failed to delete personal bento: %v\n", err)
		return c.String(http.StatusInternalServerError, "Failed to delete perosnal bento.")
	}

	return c.String(http.StatusOK, "Deleted")
}

const (
	UPDATE_ACTION_ADD    = "add"
	UPDATE_ACTION_DELETE = "delete"
	UPDATE_ACTION_UPDATE = "update"
)

type UpdatePersonalBentoAction struct {
	Method  string `json:"method" validate:"required,oneof=add update delete"`
	Name    string `json:"name" validate:"required,min=1"`
	Content string `json:"content" validate:"required"`
}

type UpdatePersonalBentoReqBody struct {
	Actions []UpdatePersonalBentoAction `json:"keyvals" validate:"required,dive"`
}

func handleUpdatePersonalBento(c echo.Context) error {
	id := c.Param("id")
	if !utils.IsValidUUIDV4(id) {
		utils.Logger().Errorf("Invalid uuid when deleting personal bento: %s\n", id)
		return c.String(http.StatusBadRequest, "Invalid uuid")
	}

	reqBody := new(UpdatePersonalBentoReqBody)
	if err := c.Bind(reqBody); err != nil {
		utils.Logger().Errorf("Failed to bind request body: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad request")
	}
	if err := c.Validate(reqBody); err != nil {
		utils.Logger().Errorf("Requests body validation failed: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad request")
	}

	// TODO: add check to verify if user has permission to perform update actions

	// make sure that the bento actually exists
	bento, err := bentomodel.GetPersonalBento(id)
	if err != nil {
		utils.Logger().Errorf("Failed to get personal bento: %v\n", err)
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "Personal bento not found.")
		}
		return c.String(http.StatusInternalServerError, "Failed to get personal bento.")
	}

	// we still need to authenticate this
	hashed := c.Request().Header.Get("X-Bento-Hashed")
	signature := c.Request().Header.Get("X-Bento-Signature")

	decodedHashed, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		utils.Logger().Errorf("Failed to decode base64 hashed challenge: %s\n", err)
		return c.String(http.StatusInternalServerError, "Failed to decode hashed challenge")
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		utils.Logger().Errorf("Failed to decode base64 signature: %s\n", err)
		return c.String(http.StatusInternalServerError, "Failed to decode signature")
	}

	err = service.VerifyBentoSignature(decodedHashed, decodedSignature, []byte(bento.PubKey))
	if err != nil {
		utils.Logger().Errorf("Failed to verify bento signature: %v\n", err)
		return c.String(http.StatusUnauthorized, "Invalid signature")
	}

	// begin transaction to allow all operations to commit at the same time
	tx, err := database.DB().Begin()
	if err != nil {
		utils.Logger().Errorf("Failed to begin transaction: %v\n", err)
		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	// TODO: add history record for tracing, this will be a version like implementation
	for _, action := range reqBody.Actions {
		switch action.Method {
		case UPDATE_ACTION_ADD:
			_, err = tx.Exec("INSERT INTO personal_bento_entries (name, content, personal_bento_id) VALUES ($1, $2, $3);", action.Name, action.Content, bento.Id)
			if err != nil {
				utils.Logger().Errorf("Failed to perform 'add' action: %v\n", err)
				if err := tx.Rollback(); err != nil {
					utils.Logger().Errorf("Failed to rollback: %v\n", err)
				}
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		case UPDATE_ACTION_DELETE:
			_, err = tx.Exec("DELETE FROM personal_bento_entries WHERE name = $1 AND personal_bento_id = $2;", action.Name, bento.Id)
			if err != nil {
				utils.Logger().Errorf("Failed to perform 'delete' action: %v\n", err)
				if err := tx.Rollback(); err != nil {
					utils.Logger().Errorf("Failed to rollback: %v\n", err)
				}
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		case UPDATE_ACTION_UPDATE:
			_, err = tx.Exec("UPDATE personal_bento_entries SET name = $1, content = $2 WHERE name = $1 AND personal_bento_id = $3", action.Name, action.Content, bento.Id)
			if err != nil {
				utils.Logger().Errorf("Failed to perform 'update' action: %v\n", err)
				if err := tx.Rollback(); err != nil {
					utils.Logger().Errorf("Failed to rollback: %v\n", err)
				}
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		utils.Logger().Errorf("Failed to commit actions: %v\n", err)
		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return c.String(http.StatusOK, "Updated")
}
