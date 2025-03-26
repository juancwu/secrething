package auth

import (
	"context"
	"time"

	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/juancwu/konbini/server/api/errors"
	"github.com/juancwu/konbini/server/api/helpers"
	"github.com/juancwu/konbini/server/db"
	"github.com/labstack/echo/v4"
)

type RegisterQueries struct {
	WebSafe bool
}

func getRegisterQueries(ctx echo.Context) RegisterQueries {
	qp := helpers.NewQueryParser(ctx)
	return RegisterQueries{
		WebSafe: qp.Bool("websafe", false),
	}
}

type RegisterRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"password"`
	Name     *string `json:"name" validate:"omitempty,omitnil,min=3,max=64,printascii"`
}

func (h *AuthHandler) Register() echo.HandlerFunc {
	// create new messages for request body validation
	messages := validator.NewValidationMessages()
	messages.SetMessage("email", "required", "Email is required")
	messages.SetMessage("email", "email", "Please enter a valid email")
	// skip password tag since it is already defined
	messages.SetMessage("name", "min", "Name must be at least {2} characters long")
	messages.SetMessage("name", "max", "Name must be at most {2} characters long")
	messages.SetMessage("name", "printascii", "Name must only consist of printable ascii characters")

	val := h.config.Validator.UseMessages(messages)

	return func(echoCtx echo.Context) error {
		processCtx, cancel := context.WithTimeout(echoCtx.Request().Context(), time.Minute)
		defer cancel()

		// Bind request body
		var reqBody RegisterRequest
		if err := echoCtx.Bind(&reqBody); err != nil {
			return err
		}
		if err := val.Validate(&reqBody); err != nil {
			return err
		}

		// qp := getRegisterQueries(echoCtx)

		conn, err := h.db.Connect()
		if err != nil {
			return errors.NewDatabaseError(err, "")
		}

		q := db.New(conn)

		// Check if user exists or not
		_, err = q.UserExistsWithEmail(processCtx, "")
		if err != nil && !db.IsNoRows(err) {
			return err
		}

		return nil
	}
}
