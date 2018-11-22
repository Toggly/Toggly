package cache

// DataCache type
type DataCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
	Flush(scopes ...string) error
}
