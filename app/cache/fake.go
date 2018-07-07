package cache

// NewFakeCache returns fake cache implementation
func NewFakeCache() (DataCache, error) {
	return &fakeCache{}, nil
}

type fakeCache struct {
}

func (c *fakeCache) Get(key string, fn func() (interface{}, error)) (data interface{}, err error) {
	return fn()
}

func (c *fakeCache) Flush(scopes ...string) {
}
