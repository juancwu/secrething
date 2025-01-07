package memcache

/*
   This file is dedicated to define prefixes for each different type
   of cache items that are used in this server. This makes a centralized
   place where all packages can get access to.
*/

const (
	EmailTokenCacheKeyPrefix string = "email_token_"
	JwtCacheKeyPrefix        string = "jwt_"
	UserCacheKeyPrefix       string = "user_"
)
