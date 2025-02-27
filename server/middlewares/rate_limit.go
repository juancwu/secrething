package middlewares

import (
	"fmt"
	"github.com/juancwu/konbini/server/memcache"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// TOTPAttempts tracks failed TOTP verification attempts
type TOTPAttempts struct {
	UserID      string    `json:"user_id"`
	Attempts    int       `json:"attempts"`
	LastAttempt time.Time `json:"last_attempt"`
}

const (
	// Max TOTP attempts before enforcing a cooldown
	MaxTOTPAttempts = 5
	// TOTP cooldown duration after exceeding max attempts
	TOTPCooldownDuration = 10 * time.Minute
	// Key prefix for memcache
	TOTPAttemptsKeyPrefix = "totp_attempts_"
)

var (
	// Mutex to prevent race conditions when updating attempts
	attemptsLock = &sync.Mutex{}
)

// LimitTOTPAttempts checks if a user has exceeded the maximum allowed TOTP verification attempts
func LimitTOTPAttempts(userID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Only apply to POST requests to relevant TOTP endpoints
			if c.Request().Method != http.MethodPost {
				return next(c)
			}

			key := TOTPAttemptsKeyPrefix + userID
			cache := memcache.Cache()

			attemptsLock.Lock()
			defer attemptsLock.Unlock()

			// Get current attempts from cache
			attemptsData, found := cache.Get(key)

			var attempts TOTPAttempts
			if found {
				attempts = attemptsData.(TOTPAttempts)

				// Check if in cooldown period
				if attempts.Attempts >= MaxTOTPAttempts {
					cooldownEnd := attempts.LastAttempt.Add(TOTPCooldownDuration)
					if time.Now().Before(cooldownEnd) {
						// Still in cooldown period
						retryAfter := int(time.Until(cooldownEnd).Seconds())
						c.Response().Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))

						return echo.NewHTTPError(
							http.StatusTooManyRequests,
							fmt.Sprintf("Too many failed attempts. Try again in %d seconds.", retryAfter),
						)
					}
					// Cooldown period expired, reset attempts
					attempts.Attempts = 0
				}
			} else {
				// Initialize new attempts tracker
				attempts = TOTPAttempts{
					UserID:      userID,
					Attempts:    0,
					LastAttempt: time.Now(),
				}
			}

			// Store original handler response writer
			resWriter := c.Response().Writer
			capturedWriter := &captureResponseWriter{ResponseWriter: resWriter, statusCode: 200}
			c.Response().Writer = capturedWriter

			// Call next handler
			err := next(c)

			// Check if TOTP verification failed based on status code
			if capturedWriter.statusCode >= 400 && capturedWriter.statusCode < 500 {
				// Increment failed attempts
				attempts.Attempts++
				attempts.LastAttempt = time.Now()

				// Calculate appropriate expiration time
				var expiration time.Duration
				if attempts.Attempts >= MaxTOTPAttempts {
					expiration = TOTPCooldownDuration
					log.Warn().
						Str("user_id", userID).
						Int("attempts", attempts.Attempts).
						Str("ip", c.RealIP()).
						Str("user_agent", c.Request().UserAgent()).
						Msg("SECURITY: User exceeded maximum TOTP attempts")
				} else {
					// Regular expiration
					expiration = TOTPCooldownDuration
				}

				// Store updated attempts in cache
				cache.Set(key, attempts, expiration)

				// Add X-RateLimit headers to response
				remainingAttempts := MaxTOTPAttempts - attempts.Attempts
				if remainingAttempts < 0 {
					remainingAttempts = 0
				}
				c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", MaxTOTPAttempts))
				c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remainingAttempts))
				c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", attempts.LastAttempt.Add(TOTPCooldownDuration).Unix()))
			} else if err == nil && capturedWriter.statusCode < 300 {
				// Successful verification - reset the counter
				cache.Delete(key)
			}

			return err
		}
	}
}

// ResetTOTPAttempts resets the TOTP attempt counter for a user
func ResetTOTPAttempts(userID string) {
	key := TOTPAttemptsKeyPrefix + userID
	memcache.Cache().Delete(key)
}

// RecordFailedTOTPAttempt increments the failed TOTP attempt counter for a user
func RecordFailedTOTPAttempt(userID string) {
	key := TOTPAttemptsKeyPrefix + userID
	cache := memcache.Cache()

	attemptsLock.Lock()
	defer attemptsLock.Unlock()

	var attempts TOTPAttempts
	attemptsData, found := cache.Get(key)

	if found {
		attempts = attemptsData.(TOTPAttempts)
	} else {
		attempts = TOTPAttempts{
			UserID:      userID,
			Attempts:    0,
			LastAttempt: time.Now(),
		}
	}

	attempts.Attempts++
	attempts.LastAttempt = time.Now()

	// Set appropriate expiration
	expiration := TOTPCooldownDuration
	cache.Set(key, attempts, expiration)

	if attempts.Attempts >= MaxTOTPAttempts {
		log.Warn().
			Str("user_id", userID).
			Int("attempts", attempts.Attempts).
			Msg("SECURITY: User exceeded maximum TOTP attempts")
	}
}

// captureResponseWriter wraps http.ResponseWriter to capture status code
type captureResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (crw *captureResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

