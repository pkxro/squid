package cache

import "github.com/pkxro/squid/internal/cache/redis"

// Cache is an interface for a cache client
type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Exists(key string) (bool, error)
}

// Manager is a wrapper for a redis client
// It provides a simple interface to set and get values from redis
type Manager struct {
	Client *redis.Cache
}

// NewCacheManager returns a new instance of Cacher
func NewCacheManager(r *redis.Cache) *Manager {
	return &Manager{
		Client: r,
	}
}
