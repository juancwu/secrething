package auth

import (
	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/juancwu/konbini/internal/server/middleware"
	"github.com/labstack/echo/v4"
)

func Configure(g *echo.Group, v *validator.Validator) {
	g.POST(
		"/sign-up",
		createUser,
		middleware.SetValidator(v, createUserBody{}),
	)

	g.POST(
		"/sign-in",
		signinUser,
		middleware.SetValidator(v, signinBody{}),
	)
}
