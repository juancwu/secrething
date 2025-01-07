package middlewares

import (
	"errors"
	"konbini/server/db"
	"konbini/server/memcache"
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
	AllowTokens []string
	ExtractFrom uint32
	ExtractUser bool
	Connector   *db.DBConnector
}

// ProtectFull is a shortcut function to protect a route which only allow full tokens.
func ProtectFull(connector *db.DBConnector) echo.MiddlewareFunc {
	return ProtectWithConfig(ProtectConfig{
		AllowTokens: []string{services.FULL_USER_TOKEN_TYPE},
		ExtractFrom: EXTRACT_FROM_BEARER,
		ExtractUser: true,
		Connector:   connector,
	})
}

// ProtectPartial is a shortcut function to protect a route which only allow partial tokens.
func ProtectPartial(connector *db.DBConnector) echo.MiddlewareFunc {
	return ProtectWithConfig(ProtectConfig{
		AllowTokens: []string{services.PARTIAL_USER_TOKEN_TYPE},
		ExtractFrom: EXTRACT_FROM_BEARER,
		ExtractUser: true,
		Connector:   connector,
	})
}

// ProtectAll is a shortcut function to protect a route which only allow partial and full tokens.
func ProtectAll(connector *db.DBConnector) echo.MiddlewareFunc {
	return ProtectWithConfig(ProtectConfig{
		AllowTokens: []string{services.PARTIAL_USER_TOKEN_TYPE, services.FULL_USER_TOKEN_TYPE},
		ExtractFrom: EXTRACT_FROM_BEARER,
		ExtractUser: true,
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

			claims, err := services.ParseUnverifyJWT(token)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to parse unverified jwt")
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}
			isAllowed := false
			for _, allow := range cfg.AllowTokens {
				if allow == claims.Type {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token type not allowed.")
			}

			// verify
			j, err := services.VerifyJWTString(c.Request().Context(), token, cfg.Connector)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to verify jwt")
				return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}
			// cache the jwt item for faster access in next verification
			memcache.CacheJWT(j)
			c.Set("jwt_id", j.ID)

			if cfg.ExtractUser {
				user, err := memcache.GetUser(j.UserID)
				if err != nil {
					logger.Warn().Err(err).Msg("Failed to get user from memory cache. Trying to get from database.")
					conn, err := cfg.Connector.Connect()
					if err != nil {
						return err
					}
					q := db.New(conn)
					dbUser, err := q.GetUserById(c.Request().Context(), j.UserID)
					if err != nil {
						conn.Close()
						return err
					}

					user = &dbUser
					conn.Close()
				}

				// renew cache for the user since it has been accessed
				memcache.CacheUser(user)

				c.Set("user_id", user.ID)
			}

			return next(c)
		}
	}
}

func GetJWT(c echo.Context) (*db.Jwt, error) {
	id, ok := c.Get("jwt_id").(string)
	if !ok {
		return nil, ErrNoJwtFound
	}
	return memcache.GetJWT(id)
}

func GetUser(c echo.Context) (*db.User, error) {
	id, ok := c.Get("user_id").(string)
	if !ok {
		return nil, ErrNoUserFound
	}
	return memcache.GetUser(id)
}
