package cachedapi

import (
	"encoding/json"
	"log"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/pkg/cache"
)

// NewCachedAPI returns cached API implementation
func NewCachedAPI(engine api.TogglyAPI, cache cache.DataCache) api.TogglyAPI {
	return &cachedAPI{engine: engine, cache: cache}
}

func withCache(cache cache.DataCache, key string, fn func() (interface{}, error)) ([]byte, error) {
	bytes, err := cache.Get(key)
	if err != nil {
		return nil, err
	}
	if bytes != nil {
		log.Printf("[DEBUG] From cache: %v", key)
		return bytes, nil
	}
	data, err := fn()
	if err != nil {
		return nil, err
	}
	bytes, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = cache.Set(key, bytes)
	if err != nil {
		log.Printf("[ERROR] Can't save data to cache: %v", err)
		return nil, err
	}
	return bytes, nil
}

type cachedAPI struct {
	engine api.TogglyAPI
	cache  cache.DataCache
}

func (c *cachedAPI) ForOwner(owner string) api.OwnerAPI {
	if c.cache == nil {
		return c.engine.ForOwner(owner)
	}
	return &cachedOwnerAPI{
		owner:  owner,
		engine: c.engine.ForOwner(owner),
		cache:  c.cache,
	}
}

type cachedOwnerAPI struct {
	owner  string
	engine api.OwnerAPI
	cache  cache.DataCache
}

func (c *cachedOwnerAPI) Projects() api.ProjectAPI {
	return &cachedProjectAPI{
		owner:  c.owner,
		engine: c.engine.Projects(),
		cache:  c.cache,
	}
}
