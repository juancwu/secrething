package memcache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var cache *gocache.Cache

func Cache() *gocache.Cache {
	if cache != nil {
		return cache
	}

	cache = gocache.New(time.Minute*5, time.Minute*10)

	return cache
}
