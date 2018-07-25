package cache

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/Toggly/core/app/rest"
)

// DataCache defines cache interface
type DataCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
	Flush(scopes ...string)
	Enabled() bool
}

func getKeyFromRequest(r *http.Request) string {
	return rest.CtxOwner(r) + "::" + r.URL.String()
}

func findSplitter(data []byte, val byte) int {
	for i, v := range data {
		if v == val {
			return i
		}
	}
	return -1
}

// Cached implements http.HandlerFunc caching
func Cached(next http.HandlerFunc, cache DataCache) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if !cache.Enabled() {
			next.ServeHTTP(w, r)
			return
		}
		key := getKeyFromRequest(r)

		log.Printf("[DEBUG] Cache search for key: %s", key)

		data, err := cache.Get(key)
		if err != nil {
			log.Printf("[ERROR] Cache error: %v", err)
			data = nil
		}

		if data != nil {
			log.Printf("[DEBUG] Cache found for key: %s", key)
			decomposeAndWriteData(key, data, w)
			return
		}

		log.Printf("[DEBUG] Cache NOT found for key: %s", key)

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		heareds := make(map[string][]string)
		for k, v := range recorder.HeaderMap {
			w.Header()[k] = v
			heareds[k] = v
		}

		item := &cacheItem{
			Headers: heareds,
			Body:    recorder.Body.Bytes(),
			Code:    recorder.Code,
		}

		w.WriteHeader(recorder.Code)
		w.Write(item.Body)

		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("[ERROR] Can't marshal cached item: %v", err)
			return
		}
		cache.Set(key, itemBytes)
	}

	return http.HandlerFunc(fn)
}

type cacheItem struct {
	Headers map[string][]string `json:"headers"`
	Body    []byte              `json:"body"`
	Code    int                 `json:"code"`
}

func decomposeAndWriteData(key string, data []byte, w http.ResponseWriter) {
	var item cacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		log.Printf("[ERROR] Can't unmarshal cached item: %v", err)
		return
	}
	for k, v := range item.Headers {
		w.Header()[k] = v
	}
	w.WriteHeader(item.Code)
	w.Write(item.Body)
}
