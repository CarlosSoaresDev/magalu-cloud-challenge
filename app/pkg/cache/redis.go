package cache

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("redis: nil")
)

type CacheClient interface {
	CheckCache() bool
	Set(key string, item interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) (*int64, error)
}

type cacheClient struct {
	cache   *redis.Client
	context context.Context
}

// New creates a new instance of cacheClient with a Redis client.
// It initializes the Redis client using the address and password
// from the environment variables REDIS_HOST_ADDRESS and REDIS_HOST_PASSWORD.
// The Redis client is configured to use database 0.
// Returns a pointer to the initialized cacheClient.
func New() *cacheClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST_ADDRESS")),
		Password: os.Getenv("REDIS_HOST_PASSWORD"),
		DB:       0,
	})

	return &cacheClient{
		context: context.Background(),
		cache:   rdb,
	}
}

// CheckCache pings the Redis cache to check if it is available.
// It returns true if the cache is reachable and false otherwise.
func (c *cacheClient) CheckCache() bool {
	_, err := c.cache.Ping(c.context).Result()
	return err == nil
}

// Set stores an item in the cache with the specified key and expiration duration.
// If the operation fails, it returns an error.
//
// Parameters:
//
//	key - the key under which the item will be stored
//	item - the item to be stored in the cache
//	expiration - the duration for which the item should remain in the cache
//
// Returns:
//
//	  error - an error if the operation fails, otherwise nil
//		*int64 - A pointer to the number of keys that were removed.
//		error - An error if the delete operation fails.
func (c *cacheClient) Set(key string, item interface{}, expiration time.Duration) error {
	return c.cache.Set(c.context, key, item, expiration).Err()
}

// Get retrieves the value associated with the given key from the cache.
// It returns the value as a byte slice and an error if the operation fails.
//
// Parameters:
//
//	key - The key for which the value needs to be retrieved.
//
// Returns:
//
//	[]byte - The value associated with the key.
//	error  - An error if the retrieval fails.
func (c *cacheClient) Get(key string) ([]byte, error) {
	return c.cache.Get(c.context, key).Bytes()
}

// Delete removes the specified key from the cache.
// It returns the number of keys that were removed and an error if the operation fails.
//
// Parameters:
//
//	key - The key to be deleted from the cache.
//
// Returns:
//
//	*int64 - A pointer to the number of keys that were removed.
//	error - An error if the delete operation fails.
func (c *cacheClient) Delete(key string) (*int64, error) {
	result, err := c.cache.Del(c.context, key).Result()

	if err != nil {
		return nil, err
	}

	return &result, nil
}
