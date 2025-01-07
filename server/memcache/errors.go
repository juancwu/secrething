package memcache

import "errors"

var (
	ErrNotFound error = errors.New("Not found in memory cache")
)
