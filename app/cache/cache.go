package cache

// DataCache defines cache interface
type DataCache interface {
	Get(key string, fn func() ([]byte, error)) (data []byte, err error)
	Flush(scopes ...string)
}
