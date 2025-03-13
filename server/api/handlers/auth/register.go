package auth

import (
	"github.com/juancwu/konbini/server/api/helpers"
	"github.com/juancwu/konbini/server/errors"
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
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=32,password"`
	Name     string `json:"name" validate:"required,min=3,max=32"`
}

func (h *AuthHandler) Register(ctx echo.Context) error {
	// qp := getRegisterQueries(ctx)

	_, err := h.db.Connect()
	if err != nil {
		return errors.NewDatabaseError(err, "")
	}

	return nil
}
