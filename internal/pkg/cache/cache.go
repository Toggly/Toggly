package cache

// DataCache defines cache interface
type DataCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
	Flush(scopes ...string)
	Enabled() bool
}
