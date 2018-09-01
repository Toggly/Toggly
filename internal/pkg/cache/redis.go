package cache

// NewRedisCache returns Redis cache implementation
func NewRedisCache() (DataCache, error) {
	return &redisCache{}, nil
}

type redisCache struct {
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
