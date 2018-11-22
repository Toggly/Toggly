package cache

import (
	"fmt"
)

// NewInMemoryCache returns in-memory cache implementation
func NewInMemoryCache() DataCache {
	return &InMemoryCache{
		Storage: make(map[string][]byte, 0),
	}
}

// InMemoryCache type
type InMemoryCache struct {
	Storage map[string][]byte
}

// Get cached data by key
func (c *InMemoryCache) Get(key string) (data []byte, err error) {
	fmt.Printf("[DEBUG] Cache get key: %s\n", key)
	return c.Storage[key], nil
}

// Set cache for key
func (c *InMemoryCache) Set(key string, data []byte) error {
	fmt.Printf("[DEBUG] Cache set key: %s\n", key)
	c.Storage[key] = data
	return nil
}

// Flush data
func (c *InMemoryCache) Flush(scopes ...string) error {
	for _, s := range scopes {
		fmt.Printf("[DEBUG] Invalidate cache for key: %s\n", s)
		delete(c.Storage, s)
	}
	return nil
}
