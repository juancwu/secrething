package handlers

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"konbini/server/config"
	"konbini/server/db"
	"konbini/server/middlewares"
	"konbini/server/services"
	"konbini/server/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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

		err = q.AddUserToGroup(
			ctx,
			db.AddUserToGroupParams{
				UserID:    user.ID,
				GroupID:   groupId,
				CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
			},
		)
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

// DeleteGroup deletes a group if the requesting user is the owner of the group
func DeleteGroup(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		groupId := c.Param("id")
		if groupId == "" {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Missing group id",
			}
		}
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}

		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		q := db.New(conn)

		exists, err := q.ExistsGroupWithIdOwnedByUser(c.Request().Context(), db.ExistsGroupWithIdOwnedByUserParams{
			ID:      groupId,
			OwnerID: user.ID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "Group not found",
					InternalError: err,
				}
			}
			return err
		}
		if exists != 1 {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Group not found",
				InternalError: err,
			}
		}

		err = q.RemoveGroupByID(c.Request().Context(), groupId)
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

type InviteUsersToJoinGroupRequest struct {
	GroupID string   `json:"group_id" validate:"required,uuid4"`
	Emails  []string `json:"emails" validate:"gt=0,dive,email"`
}

func InviteUsersToJoinGroup(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return err
		}
		body, err := middlewares.GetJsonBody[InviteUsersToJoinGroupRequest](c)
		if err != nil {
			return err
		}
		cfg, err := config.Global()
		if err != nil {
			return err
		}
		logger := middlewares.GetLogger(c)

		// check if user is owner of group
		conn, err := connector.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		q := db.New(conn)

		group, err := q.GetGroupByIDOwendByUser(
			c.Request().Context(),
			db.GetGroupByIDOwendByUserParams{
				ID:      body.GroupID,
				OwnerID: user.ID,
			},
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return APIError{
					Code:          http.StatusBadRequest,
					PublicMessage: "No group found",
					InternalError: err,
				}
			}
			return err
		}

		// the token needs to have the
		// - invitation id
		// - expiry time
		// default expiry time is 24 hours

		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		q = q.WithTx(tx)

		params := services.SendGroupInvitationEmailsParams{
			InvitorName: user.Nickname,
			GroupName:   group.Name,
			Users: make([]struct {
				Name  string
				Token string
				Email string
			}, len(body.Emails)),
		}
		for i, email := range body.Emails {
			invitedUser, err := q.GetUserByEmail(c.Request().Context(), email)
			if err != nil {
				if err == sql.ErrNoRows {
					return APIError{
						Code:          http.StatusBadRequest,
						PublicMessage: fmt.Sprintf("No user with email: '%s'", email),
						InternalError: err,
					}
				}
				return err
			}
			now := time.Now()
			exp := now.Add(time.Hour * 24)
			invitationId, err := q.NewGroupInvitation(
				c.Request().Context(),
				db.NewGroupInvitationParams{
					UserID:    invitedUser.ID,
					GroupID:   body.GroupID,
					CreatedAt: utils.FormatRFC3339NanoFixed(now),
					ExpiresAt: utils.FormatRFC3339NanoFixed(exp),
				},
			)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					logger.Error().Err(err).Msg("Failed to rollback")
				}
				return err
			}

			// 36 bytes from invitation id + 30 from time
			token := make([]byte, 66)
			copy(token[:], []byte(invitationId))
			copy(token[36:], []byte(utils.FormatRFC3339NanoFixed(exp)))

			token, err = utils.EncryptAES(token, cfg.GetAesKey())
			if err != nil {
				if err := tx.Rollback(); err != nil {
					logger.Error().Err(err).Msg("Failed to rollback")
				}
				return err
			}

			// add encrypted token to list
			params.Users[i].Token = base64.URLEncoding.EncodeToString(token)
			params.Users[i].Name = invitedUser.Nickname
			params.Users[i].Email = invitedUser.Email
		}

		// commit changes
		err = tx.Commit()
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logger.Error().Err(err).Msg("Failed to rollback")
			}
			return err
		}

		// send the emails
		go func(params services.SendGroupInvitationEmailsParams) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			res, err := services.SendGroupInvitationEmails(ctx, params)
			if err != nil {
				log.Error().Err(err).Msg("Failed to send group invitations")
				return
			}
			ids := make([]string, len(res.Data))
			for i, d := range res.Data {
				ids[i] = d.Id
			}
			log.Info().Strs("email_ids", ids).Msg("Successfully sent group invitations")
		}(params)

		return c.NoContent(http.StatusCreated)
	}
}

// AcceptGroupInvitation accepts an invitation to join a group
func AcceptGroupInvitation(connector *db.DBConnector) echo.HandlerFunc {
	return func(c echo.Context) error {
		b64Token := c.QueryParam("token")
		if b64Token == "" {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Missing token",
			}
		}

		cfg, err := config.Global()
		if err != nil {
			return err
		}

		// validate the token
		token, err := base64.URLEncoding.DecodeString(b64Token)
		if err != nil {
			return APIError{
				Code:           http.StatusBadRequest,
				PublicMessage:  "Invalid token",
				PrivateMessage: "Failed to decode base64 token",
				InternalError:  err,
			}
		}

		// decrypt the token
		token, err = utils.DecryptAES(token, cfg.GetAesKey())
		if err != nil {
			return APIError{
				Code:           http.StatusBadRequest,
				PublicMessage:  "Invalid token",
				PrivateMessage: "Failed to decrypt token",
				InternalError:  err,
			}
		}

		// 36 bytes invitation id + 30 from time
		if len(token) != 66 {
			return APIError{
				Code:           http.StatusBadRequest,
				PublicMessage:  "Invalid token",
				PrivateMessage: "Token is not 66 bytes long",
			}
		}

		invitationId := token[:36]
		exp := token[36:]

		// parse expiration and check
		expT, err := time.Parse(time.RFC3339Nano, string(exp))
		if err != nil {
			return err
		}

		if time.Now().After(expT) {
			return APIError{
				Code:          http.StatusBadRequest,
				PublicMessage: "Invitation Expired",
			}
		}

		// get invitation from database
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

		invitation, err := q.GetGroupInvitationByID(c.Request().Context(), string(invitationId))
		if err != nil {
			if err == sql.ErrNoRows {
				return APIError{
					Code:           http.StatusBadRequest,
					PublicMessage:  "Invalid Invitation",
					PrivateMessage: "No invitation found in the database",
					InternalError:  err,
				}
			}
			return err
		}

		err = q.AddUserToGroup(
			c.Request().Context(),
			db.AddUserToGroupParams{
				UserID:    invitation.UserID,
				GroupID:   invitation.GroupID,
				CreatedAt: utils.FormatRFC3339NanoFixed(time.Now()),
			},
		)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = q.RemoveGroupInvitationByID(
			c.Request().Context(),
			invitation.ID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}
