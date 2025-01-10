package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"konbini/server/db"
	"konbini/server/middlewares"
	"konbini/server/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type NewBentoRequest struct {
	Name        string `json:"name" validate:"required,min=3,printascii"`
	Ingridients []struct {
		Name  string `json:"name" validate:"required,min=1,printascii"`
		Value string `json:"value"`
	} `json:"ingridients,omitempty" validate:"omitnil,omitempty,dive"`
}

func NewBento(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}
		if !user.EmailVerified {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Email must be verified before creating a new bento",
			}
		}

		body, err := middlewares.GetJsonBody[NewBentoRequest](c)
		if err != nil {
			return err
		}

		// shouldn't take more than 5 seconds to run
		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
		defer cancel()

		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()
		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		q := db.New(tx)

		exists, err := q.ExistsBentoWithNameOwnedByUser(
			ctx,
			db.ExistsBentoWithNameOwnedByUserParams{
				Name:   body.Name,
				UserID: user.ID,
			},
		)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if exists == 1 {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: fmt.Sprintf("Bento with name %s already exists.", body.Name),
			}
		}

		// create bento
		bentoID, err := q.NewBento(
			ctx,
			db.NewBentoParams{
				Name:      body.Name,
				UserID:    user.ID,
				CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
				UpdatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
			},
		)
		if err != nil {
			tx.Rollback()
			return err
		}

		if body.Ingridients != nil && len(body.Ingridients) > 0 {
			for _, ing := range body.Ingridients {
				err = q.AddIngridientToBento(
					ctx,
					db.AddIngridientToBentoParams{
						BentoID:   bentoID,
						Name:      ing.Name,
						Value:     []byte(ing.Value),
						CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
						UpdatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
					},
				)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}

		return c.JSON(http.StatusCreated, map[string]string{"bento_id": bentoID})
	}
}
