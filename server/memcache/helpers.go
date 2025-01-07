package memcache

import (
	"fmt"
	"konbini/server/services"
	"time"
)

/*
   This file is dedicated to define prefixes for each different type of cache items and helper functions
   that are used in this server. This makes a centralized place where all packages can get access to.
*/

const (
	EmailTokenCacheKeyPrefix string = "email_token_"
	JwtCacheKeyPrefix        string = "jwt_"
)

// StoreEmailTokenInCache is a helper function that stores the given email token in memory cache.
// The function uses the emailTokenCacheKeyPrefix to prefix the key before storing it.
// The function returns an error if the key already exists in cache.
func StoreEmailTokenInCache(token *services.EmailToken) error {
	return cache.Add(EmailTokenCacheKeyPrefix+token.Id, token, time.Minute*10)
}

// GetEmailTokenFromCache is a helper function that retrieves the token email with the given id.
// The function returns an error if the key does not exists, expired cache or stored value is not
// a valid service.EmailToken struct.
func GetEmailTokenFromCache(id string) (*services.EmailToken, error) {
	k, exp, found := cache.GetWithExpiration(EmailTokenCacheKeyPrefix + id)
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

// StoreJwtInCache stores the JWT in cache using a prefix "jwt_" + id + id
// The JWT is cached for 1 hour.
func StoreJwtInCache(id string, jwt *services.JWT) error {
	return cache.Add(JwtCacheKeyPrefix+id, jwt, time.Hour)
}

// GetJwtFromCache retrieves the JWT stored in memory cache.
// This function will also check if the JWT is about to expired (10 minutes)
// and renews the expiration date. Allows the easy access of the JWT
// for a long continuous time.
func GetJwtFromCache(id string) (*services.JWT, error) {
	k, exp, found := cache.GetWithExpiration(JwtCacheKeyPrefix + id)
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
		cache.Replace(JwtCacheKeyPrefix+id, jwt, time.Hour)
	}
	return jwt, nil
}
