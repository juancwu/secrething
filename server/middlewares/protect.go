package middlewares

import (
	"errors"
	"konbini/server/db"
	"konbini/server/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	// EXTRACT_FROM_BEARER means that the token would be found in the Authorization header as a Bearer token.
	EXTRACT_FROM_BEARER uint32 = 0
)

var (
	ErrNoJwtFound  error = errors.New("No jwt found")
	ErrNoUserFound error = errors.New("No user found")
)

type ProtectConfig struct {
	AllowTokens []services.TokenType
	ExtractFrom uint32
	Connector   *db.DBConnector
}

// ProtectFull is a shortcut function to protect a route which only allow full tokens.
func ProtectFull(connector *db.DBConnector) echo.MiddlewareFunc {
	return ProtectWithConfig(ProtectConfig{
		AllowTokens: []services.TokenType{services.FULL_USER_TOKEN_TYPE},
		ExtractFrom: EXTRACT_FROM_BEARER,
		Connector:   connector,
	})
}

// ProtectPartial is a shortcut function to protect a route which only allow partial tokens.
func ProtectPartial(connector *db.DBConnector) echo.MiddlewareFunc {
	return ProtectWithConfig(ProtectConfig{
		AllowTokens: []services.TokenType{services.PARTIAL_USER_TOKEN_TYPE},
		ExtractFrom: EXTRACT_FROM_BEARER,
		Connector:   connector,
	})
}

// ProtectAll is a shortcut function to protect a route which only allow partial and full tokens.
func ProtectAll(connector *db.DBConnector) echo.MiddlewareFunc {
	return ProtectWithConfig(ProtectConfig{
		AllowTokens: []services.TokenType{services.PARTIAL_USER_TOKEN_TYPE, services.FULL_USER_TOKEN_TYPE},
		ExtractFrom: EXTRACT_FROM_BEARER,
		Connector:   connector,
	})
}

// ProtectWithConfig is a middleware that checks for a jwt in the request and validates it.
// The middleware also refreshes the cache for the jwt in memory if the expiry <= 10 minutes.
func ProtectWithConfig(cfg ProtectConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := GetLogger(c)
			var token string
			switch cfg.ExtractFrom {
			case EXTRACT_FROM_BEARER:
				header := c.Request().Header.Get(echo.HeaderAuthorization)
				parts := strings.Split(header, " ")
				if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format. Expecting Bearer token.")
				}
				token = parts[1]
			default:
				return errors.New("Invalid ExtractFrom value in Protect middleware configuration.")
			}

			// verify the token
			jwt, err := services.VerifyJWT(token)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to verify JWT")
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			allowed := false
			for _, t := range cfg.AllowTokens {
				if jwt.TokenType == t {
					allowed = true
					break
				}
			}
			if !allowed {
				logger.Error().Err(err).Msg("Request made with token type that is NOT allowed.")
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			// check database if token exists
			conn, err := cfg.Connector.Connect()
			if err != nil {
				return err
			}

			q := db.New(conn)

			exists, err := q.ExistsJwtById(c.Request().Context(), jwt.ID)
			if err != nil {
				conn.Close()
				logger.Error().Err(err).Msg("Failed to fetch JWT from database. Reject.")
				return err
			}
			if exists != 1 {
				conn.Close()
				logger.Error().Msg("JWT does not exists in database. Reject.")
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			user, err := q.GetUserById(c.Request().Context(), jwt.UserID)
			if err != nil {
				conn.Close()
				logger.Error().Msg("Failed to fetch user from database. Reject.")
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			// close connection after use
			conn.Close()

			c.Set("jwt", jwt)
			c.Set("user", user)

			return next(c)
		}
	}
}

func GetJWT(c echo.Context) (*services.JWT, error) {
	jwt, ok := c.Get("jwt").(*services.JWT)
	if !ok {
		return nil, ErrNoJwtFound
	}
	return jwt, nil
}

func GetUser(c echo.Context) (db.User, error) {
	user, ok := c.Get("user").(db.User)
	if !ok {
		return db.User{}, ErrNoUserFound
	}
	return user, nil
}
