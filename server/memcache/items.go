package memcache

import (
	"konbini/server/db"
	"time"
)

func GetJWT(id string) (*db.Jwt, error) {
	k, found := cache.Get(JwtCacheKeyPrefix + id)
	if !found {
		return nil, ErrNotFound
	}
	if k, ok := k.(*db.Jwt); ok {
		return k, nil
	}
	return nil, ErrInvalidItem
}

func CacheJWT(item *db.Jwt) {
	cache.Set(JwtCacheKeyPrefix+item.ID, item, time.Hour)
}

func GetUser(id string) (*db.User, error) {
	k, found := cache.Get(UserCacheKeyPrefix + id)
	if !found {
		return nil, ErrNotFound
	}
	if k, ok := k.(*db.User); ok {
		return k, nil
	}
	return nil, ErrInvalidItem
}

func CacheUser(item *db.User) {
	cache.Set(UserCacheKeyPrefix+item.ID, item, time.Hour)
}
