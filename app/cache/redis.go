package cache

// NewRedisCache returns Redis cache implementation
func NewRedisCache() (DataCache, error) {
	return &redisCache{}, nil
}

type redisCache struct {
}

func (c *redisCache) Get(key string, fn func() ([]byte, error)) (data []byte, err error) {
	return nil, nil
}

func (c *redisCache) Flush(scopes ...string) {
}
