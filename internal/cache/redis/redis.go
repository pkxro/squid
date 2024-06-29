package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/redis.v5"
)

// IdempotencyLockKey is the key used to store the idempotency lock in redis
const IdempotencyLockKey string = "il_idempotency_lock_%s"

// Cache is a wrapper for a redis client
type Cache struct {
	db  *redis.Client
	ttl time.Duration
}

// NewRedisCache returns a new instance of RedisCache
func NewRedisCache(uri string, password string, ttl time.Duration) (*Cache, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: password,
	})

	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Cache{
		db:  db,
		ttl: ttl,
	}, nil
}

func createIdempotencyKey(attr string) string {
	return fmt.Sprintf(IdempotencyLockKey, attr)
}

// Exists checks if a key exists in Redis
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	cacheKey := createIdempotencyKey(key)

	ok, err := c.db.Exists(cacheKey).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

// Get returns a value from redis
func (c *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	var lock interface{}

	cacheKey := createIdempotencyKey(key)

	b, err := c.db.Get(cacheKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &lock); err != nil {
		return nil, err
	}

	return &lock, nil
}

// Set sets a value in redis
func (c *Cache) Set(ctx context.Context, key string, lock interface{}, exp time.Duration) error {
	cacheKey := createIdempotencyKey(key)

	b, err := json.Marshal(lock)
	if err != nil {
		return err
	}

	set := c.db.Set(cacheKey, b, exp)
	if set.Err() != nil {
		return err
	}
	return err
}