package main

import "github.com/Toggly/core/internal/pkg/cache"

func main() {}

// GetCache returns in-memory DataCache implementation
func GetCache() cache.DataCache {
	return &cache.InMemoryCache{
		Storage: make(map[string][]byte, 0),
	}
}
