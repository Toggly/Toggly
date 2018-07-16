package cache

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

// DataCache defines cache interface
type DataCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
	Flush(scopes ...string)
	Enabled() bool
}

func getKeyFromRequest(r *http.Request) string {
	return r.URL.String()
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

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			sp := findSplitter(data, 0)
			if sp > -1 {
				h := make(map[string][]string)
				if err = json.Unmarshal(data[:sp], &h); err != nil {
					log.Printf("[ERROR] %v", err)
				}
				for k, v := range h {
					w.Header()[k] = v
				}
				data = data[sp+1:]
			}
			w.Write(data)
			return
		}

		log.Printf("[DEBUG] Cache NOT found for key: %s", key)
		c := httptest.NewRecorder()
		next.ServeHTTP(c, r)

		heareds := make(map[string][]string)
		for k, v := range c.HeaderMap {
			w.Header()[k] = v
			heareds[k] = v
		}
		w.WriteHeader(c.Code)

		body := c.Body.Bytes()
		w.Write(body)

		h, _ := json.Marshal(&heareds)
		d := make([]byte, 0)
		d = append(d, h...)
		d = append(d, 0)
		d = append(d, body...)

		cache.Set(key, d)
	})

}
