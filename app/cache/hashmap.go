package cache

// NewHashMapCache returns in-memory hashmap cache implementation. For development purposes
func NewHashMapCache(enabled bool) (DataCache, error) {
	return &fakeCache{
		enabled: enabled,
		storage: make(map[string][]byte, 0),
	}, nil
}

type fakeCache struct {
	storage map[string][]byte
	enabled bool
}

func (c *fakeCache) Enabled() bool {
	return c.enabled
}

func (c *fakeCache) Get(key string) (data []byte, err error) {
	return c.storage[key], nil
}

func (c *fakeCache) Set(key string, data []byte) error {
	c.storage[key] = data
	return nil
}

func (c *fakeCache) Flush(scopes ...string) {
	// TODO implement
}
