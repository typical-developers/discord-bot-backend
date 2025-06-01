package redisutil

import "time"

type CacheOpts struct {
	Expiry *time.Duration
}
