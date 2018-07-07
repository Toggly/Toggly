package cache

// DataCache defines cache interface
type DataCache interface {
	Get(key string, fn func() (interface{}, error)) (data interface{}, err error)
	Flush(scopes ...string)
}
