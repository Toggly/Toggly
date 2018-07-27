package cache

// NewRedisCache returns Redis cache implementation
func NewRedisCache(enabled bool) (DataCache, error) {
	return &redisCache{
		enabled: enabled,
	}, nil
}

type redisCache struct {
	enabled bool
}

func (c *redisCache) Enabled() bool {
	return c.enabled
}

func (c *redisCache) Get(key string) (data []byte, err error) {
	return nil, nil
}

func (c *redisCache) Set(key string, data []byte) error {
	return nil
}

func (c *redisCache) Flush(scopes ...string) {
	// TODO implement
}
