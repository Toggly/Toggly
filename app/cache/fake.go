package cache

// NewFakeCache returns fake cache implementation
func NewFakeCache() (DataCache, error) {
	return &fakeCache{}, nil
}

type fakeCache struct {
}

func (c *fakeCache) Get(key string, fn func() ([]byte, error)) (data []byte, err error) {
	if data, err = fn(); err != nil {
		return data, err
	}
	return data, nil
}

func (c *fakeCache) Flush(scopes ...string) {
}
