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

func joinBytes(data1 []byte, data2 []byte, splitter byte) []byte {
	d := make([]byte, 0)
	d = append(d, data1...)
	d = append(d, splitter)
	d = append(d, data2...)
	return d
}

// Cached implements http.HandlerFunc caching
func Cached(fn func(wr http.ResponseWriter, req *http.Request), cache DataCache) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
		fn(c, r)

		heareds := make(map[string][]string)
		for k, v := range c.HeaderMap {
			w.Header()[k] = v
			heareds[k] = v
		}
		w.WriteHeader(c.Code)

		body := c.Body.Bytes()
		w.Write(body)

		h, _ := json.Marshal(&heareds)
		d := joinBytes(h, body, 0)

		cache.Set(key, d)
	})
}
