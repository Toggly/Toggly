package cache

import (
	"fmt"
)

// NewHashMapCache returns in-memory hashmap cache implementation. For development purposes
func NewHashMapCache() (DataCache, error) {
	return &fakeCache{
		storage: make(map[string][]byte, 0),
	}, nil
}

type fakeCache struct {
	storage map[string][]byte
}

func (c *fakeCache) Get(key string) (data []byte, err error) {
	fmt.Printf("[DEBUG] Cache get key: %s\n", key)
	return c.storage[key], nil
}

func (c *fakeCache) Set(key string, data []byte) error {
	fmt.Printf("[DEBUG] Cache set key: %s\n", key)
	c.storage[key] = data
	return nil
}

func (c *fakeCache) Flush(scopes ...string) {
	for _, s := range scopes {
		fmt.Printf("[DEBUG] Invalidate cache for key: %s\n", s)
		delete(c.storage, s)
	}
}
