package router

import (
	"database/sql"
	"encoding/base64"
	"net/http"

	"github.com/juancwu/konbini/server/middleware"
	bentomodel "github.com/juancwu/konbini/server/models/bento"
	usermodel "github.com/juancwu/konbini/server/models/user"
	"github.com/juancwu/konbini/server/service"
	"github.com/juancwu/konbini/server/utils"
	"github.com/labstack/echo/v4"
)

func SetupBentoRoutes(e *echo.Echo) {
	e.POST("/bento/personal/new", handleNewPersonalBento, middleware.JwtAuthMiddleware)
	e.GET("/bento/personal/:id", handleGetPersonalBento)
	e.GET("/bento/personal/list", handleListPersonalBentos, middleware.JwtAuthMiddleware)
}

type NewPersonalBentoReqBody struct {
	Name       string `json:"name" validate:"required,min=1"`
	PublickKey string `json:"public_key" validate:"required,min=1"`
	Content    string `json:"content" validate:"required"`
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

	bentoId, err := bentomodel.NewPersonalBento(claims.UserId, reqBody.Name, reqBody.PublickKey, reqBody.Content)
	if err != nil {
		utils.Logger().Errorf("Failed to create new personal bento: %v\n", err)
		return c.String(http.StatusInternalServerError, "Failed to create personal bento. Please try again later.")
	}
	utils.Logger().Info("New personal bento created.", "user_id", claims.UserId, "bento_id", bentoId)

	return c.String(http.StatusCreated, bentoId)
}

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
