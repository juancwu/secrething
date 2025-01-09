package handlers

import (
	"fmt"
	"konbini/server/memcache"
	"konbini/server/services"
	"time"
)

// storeEmailTokenInCache is a helper function that stores the given email token in memory cache.
// The function uses the emailTokenCacheKeyPrefix to prefix the key before storing it.
// The function returns an error if the key already exists in cache.
func storeEmailTokenInCache(token *services.EmailToken) error {
	cache := memcache.Cache()
	return cache.Add("ve_"+token.Id, token, time.Minute*10)
}

// getEmailTokenFromCache is a helper function that retrieves the token email with the given id.
// The function returns an error if the key does not exists, expired cache or stored value is not
// a valid service.EmailToken struct.
func getEmailTokenFromCache(id string) (*services.EmailToken, error) {
	cache := memcache.Cache()
	k, exp, found := cache.GetWithExpiration("ve_" + id)
	if !found {
		return nil, memcache.ErrNotFound
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
