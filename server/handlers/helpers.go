package handlers

import (
	"fmt"
	"konbini/server/memcache"
	"konbini/server/services"
	"time"
)

// emailTokenCacheKeyPrefix is the prefix that is added when an email token is stored in cache
const emailTokenCacheKeyPrefix string = "email_token_"
const jwtCacheKeyPrefix string = "jwt_"

// storeEmailTokenInCache is a helper function that stores the given email token in memory cache.
// The function uses the emailTokenCacheKeyPrefix to prefix the key before storing it.
// The function returns an error if the key already exists in cache.
func storeEmailTokenInCache(token *services.EmailToken) error {
	cache := memcache.Cache()
	return cache.Add(emailTokenCacheKeyPrefix+token.Id, token, time.Minute*10)
}

// getEmailTokenFromCache is a helper function that retrieves the token email with the given id.
// The function returns an error if the key does not exists, expired cache or stored value is not
// a valid service.EmailToken struct.
func getEmailTokenFromCache(id string) (*services.EmailToken, error) {
	cache := memcache.Cache()
	k, exp, found := cache.GetWithExpiration(emailTokenCacheKeyPrefix + id)
	if !found {
		return nil, fmt.Errorf("No email token found with id: %s", id)
	}
	if time.Now().UTC().After(exp) {
		return nil, fmt.Errorf("Email token cache expired. ID: %s", id)
	}
	token, ok := k.(*services.EmailToken)
	if !ok {
		return nil, fmt.Errorf("Invalid email token type.")
	}
	return token, nil
}

// storeJwtInCache stores the JWT in cache using a prefix "jwt_" + id + id
// The JWT is cached for 1 hour.
func storeJwtInCache(id string, jwt *services.JWT) error {
	cache := memcache.Cache()
	return cache.Add(jwtCacheKeyPrefix+id, jwt, time.Hour)
}

// getJwtFromCache retrieves the JWT stored in memory cache.
// This function will also check if the JWT is about to expired (10 minutes)
// and renews the expiration date. Allows the easy access of the JWT
// for a long continuous time.
func getJwtFromCache(id string) (*services.JWT, error) {
	cache := memcache.Cache()
	k, exp, found := cache.GetWithExpiration(jwtCacheKeyPrefix + id)
	if !found {
		return nil, fmt.Errorf("No JWT found with id: %s", id)
	}
	now := time.Now().UTC()
	if now.After(exp) {
		return nil, fmt.Errorf("JWT cache with id %s has expired.", id)
	}
	jwt, ok := k.(*services.JWT)
	if !ok {
		return nil, fmt.Errorf("JWT cache with id %s has invalid type.", id)
	}
	diff := exp.Sub(now)
	if diff < 0 {
		diff = -diff
	}
	if diff <= 10*time.Minute {
		// update the time
		cache.Replace(jwtCacheKeyPrefix+id, jwt, time.Hour)
	}
	return jwt, nil
}
