package auth

import (
	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/juancwu/secrething/internal/server/middleware"
	"github.com/labstack/echo/v4"
)

// Configure registers all auth-related routes to the given router with validators.
func Configure(g *echo.Group, v *validator.Validator) {
	// User registration
	g.POST(
		"/sign-up",
		createUser,
		middleware.SetValidator(v, getCreateUserRequestMessages()),
	)
	
	// CLI registration
	g.POST(
		"/cli/sign-up",
		createUser,
		middleware.SetValidator(v, getCreateUserRequestMessages()),
	)

	// Web authentication
	g.POST(
		"/sign-in",
		loginUser,
		middleware.SetValidator(v, getSigninRequestMessages()),
	)

	g.POST(
		"/sign-in/totp",
		verifyTOTP,
		middleware.SetValidator(v, getTotpVerifyRequestMessages()),
	)

	// CLI authentication
	g.POST(
		"/cli/sign-in",
		cliLogin,
		middleware.SetValidator(v, getSigninRequestMessages()),
	)

	g.POST(
		"/cli/sign-in/totp",
		verifyCliTOTP,
		middleware.SetValidator(v, getCliTotpVerifyRequestMessages()),
	)

	// Token management - no validation needed for refresh/logout as they're optional
	g.POST("/refresh", refreshTokens)
	g.POST("/logout", logout)

	// TOTP management routes
	totp := g.Group("/totp", middleware.Protected())
	totp.DELETE("/remove", func(c echo.Context) error {
		return nil // TODO: Implement TOTP removal
	})
	totp.POST("/activate", func(c echo.Context) error {
		return nil // TODO: Implement TOTP activation
	})
	totp.POST("/verify", func(c echo.Context) error {
		return nil // TODO: Implement TOTP verification
	})

	// API token management - requires authentication
	tokens := g.Group("/token", middleware.Protected())
	tokens.POST(
		"/api",
		createAPIToken,
	)
	tokens.DELETE(
		"/api",
		revokeAPIToken,
		middleware.SetValidator(v, getRevokeAPITokenRequestMessages()),
	)
}
