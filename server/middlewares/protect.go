package middlewares

import (
	"database/sql"
	"errors"
	"konbini/server/db"
	"konbini/server/memcache"
	"konbini/server/services"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	// EXTRACT_FROM_BEARER means that the token would be found in the Authorization header as a Bearer token.
	EXTRACT_FROM_BEARER uint32 = 0
)

var (
	ErrNoJwtFound  error = errors.New("No authToken found")
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

// ProtectWithConfig is a middleware that checks for a authToken in the request and validates it.
// The middleware also refreshes the cache for the authToken in memory if the expiry <= 10 minutes.
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
			authToken, err := services.VerifyAuthToken(token)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to verify AuthToken")
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}

			allowed := false
			for _, t := range cfg.AllowTokens {
				if authToken.TokenType == t {
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

			// before querying, check memory cache
			_, found := memcache.Cache().Get("auth_token_" + authToken.ID)
			if !found {
				exists, err := q.ExistsAuthTokenById(c.Request().Context(), authToken.ID)
				if err != nil {
					conn.Close()
					if err == sql.ErrNoRows {
						return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
					}
					logger.Error().Err(err).Msg("Failed to fetch AuthToken from database. Reject.")
					return err
				}
				if exists != 1 {
					conn.Close()
					logger.Error().Msg("AuthToken does not exists in database. Reject.")
					return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				}
				memcache.Cache().Set("auth_token_"+authToken.ID, authToken.ID, time.Minute*10)
			}

			user, err := q.GetUserById(c.Request().Context(), authToken.UserID)
			if err != nil {
				conn.Close()
				if err == sql.ErrNoRows {
					logger.Error().Msg("No user found in database")
					return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				}
				logger.Error().Msg("Failed to fetch user from database. Reject.")
				return err
			}

			// close connection after use
			conn.Close()

			c.Set("authToken", authToken)
			c.Set("user", user)

			return next(c)
		}
	}
}

func GetJWT(c echo.Context) (*services.AuthToken, error) {
	authToken, ok := c.Get("authToken").(*services.AuthToken)
	if !ok {
		return nil, ErrNoJwtFound
	}
	return authToken, nil
}

func GetUser(c echo.Context) (db.User, error) {
	user, ok := c.Get("user").(db.User)
	if !ok {
		return db.User{}, ErrNoUserFound
	}
	return user, nil
}
