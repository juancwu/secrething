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

// NewGroupRequest is the request body to create a new group
type NewGroupRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50,printascii"`
}

// NewGroup handles request to create new groups
func NewGroup(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := middlewares.GetLogger(c)
		user, err := middlewares.GetUser(c)
		if err != nil {
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Failed to get user from context",
				InternalError:  err,
			}
		}
		body, err := middlewares.GetJsonBody[NewGroupRequest](c)
		if err != nil {
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Failed to get json body from context",
				InternalError:  err,
			}
		}

		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		q := db.New(conn)

		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Minute)
		defer cancel()

		exists, err := q.ExistsGroupOwnedByUser(ctx, db.ExistsGroupOwnedByUserParams{
			OwnerID: user.ID,
			Name:    body.Name,
		})
		if err != nil && err != sql.ErrNoRows {
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Error when checking if group exists under requesting user",
				InternalError:  err,
			}
		}
		if exists == 1 {
			return APIError{
				Code:           http.StatusBadRequest,
				PublicMessage:  fmt.Sprintf("Group with name: '%s' already exists.", body.Name),
				PrivateMessage: "Duplicate group",
			}
		}

		tx, err := conn.Begin()
		if err != nil {
			return err
		}
		q = q.WithTx(tx)

		groupId, err := q.NewGroup(ctx, db.NewGroupParams{
			Name:      body.Name,
			OwnerID:   user.ID,
			CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
			UpdatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
		})
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error().Err(err).Msg("failed to rollback")
			}
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Failed to create new group",
				InternalError:  err,
			}
		}

		err = q.AddUserToGroup(ctx, db.AddUserToGroupParams{UserID: user.ID, GroupID: groupId})
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error().Err(err).Msg("failed to rollback")
			}
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Failed add user to group",
				InternalError:  err,
			}
		}

		err = tx.Commit()
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error().Err(err).Msg("failed to rollback")
			}
			return APIError{
				Code:           http.StatusInternalServerError,
				PrivateMessage: "Failed to commit changes",
				InternalError:  err,
			}
		}

		return c.JSON(http.StatusCreated, map[string]string{"group_id": groupId})
	}
}
