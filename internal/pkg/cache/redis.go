package cache

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

// NewRedisCache returns redis cache implementation
func NewRedisCache(url string) DataCache {
	c, err := redis.DialURL(url)
	if err != nil {
		log.Fatalf("Can't connect to Redis url `%s` : %v", url, err)
	}
	return &RedisCache{Conn: c}
}

// RedisCache type
type RedisCache struct {
	Conn redis.Conn
}

// Get bytes by key
func (c *RedisCache) Get(key string) ([]byte, error) {
	fmt.Printf("[DEBUG] Cache get key: %s\n", key)
	data, err := redis.Bytes(c.Conn.Do("GET", key))
	if err == redis.ErrNil {
		return nil, nil
	}
	return data, err
}

// Set bytes by key
func (c *RedisCache) Set(key string, data []byte) error {
	fmt.Printf("[DEBUG] Cache set key: %s\n", key)
	_, err := c.Conn.Do("SET", key, data)
	return err
}

// Flush cached data
func (c *RedisCache) Flush(scopes ...string) error {
	for _, key := range scopes {
		fmt.Printf("[DEBUG] Invalidate cache for key: %s\n", key)
		if _, err := c.Conn.Do("DEL", key); err != nil {
			return err
		}

	}
	return nil
}
