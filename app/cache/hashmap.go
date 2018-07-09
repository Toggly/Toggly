package cache

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
	return c.storage[key], nil
}

func (c *fakeCache) Set(key string, data []byte) error {
	c.storage[key] = data
	return nil
}

func (c *fakeCache) Flush(scopes ...string) {
}
