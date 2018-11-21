package main

import (
	"fmt"

	in "github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/pkg/cache"
)

func main() {}

// GetCache returns in-memory DataCache implementation
func GetCache(parameters map[string]string) cache.DataCache {
	fmt.Printf("params: %v", parameters)
	return &in.InMemoryCache{
		Storage: make(map[string][]byte, 0),
	}
}
