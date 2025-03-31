package auth

import (
	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/juancwu/secrething/internal/server/middleware"
	"github.com/labstack/echo/v4"
)

func Configure(g *echo.Group, v *validator.Validator) {
	g.POST(
		"/sign-up",
		createUser,
		middleware.SetValidator(v, getCreateUserRequestMessages()),
	)

	g.POST(
		"/sign-in",
		signinUser,
		middleware.SetValidator(v, getSigninRequestMessages()),
	)

	totp := g.Group("/totp", middleware.Protected())
	totp.DELETE("/remove", func(c echo.Context) error {
		return nil
	})
	totp.POST("/activate", func(c echo.Context) error {
		return nil
	})
	totp.POST("/verifiy", func(c echo.Context) error {
		return nil
	})
}
