package main

import (
	in "github.com/Toggly/core/internal/pkg/cache"
)

func main() {}

// GetCache returns in-memory DataCache implementation
func GetCache(parameters map[string]string) interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
	Flush(scopes ...string) error
} {
	return &in.InMemoryCache{
		Storage: make(map[string][]byte, 0),
	}
}
