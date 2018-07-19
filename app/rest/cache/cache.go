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
		w.WriteHeader(recorder.Code)

		body := recorder.Body.Bytes()
		w.Write(body)

		h, err := json.Marshal(&heareds)
		if err != nil {
			log.Printf("[ERROR] Can't parse response headers: %v", err)
			return
		}

		cache.Set(key, composeData(h, body))
	}

	return http.HandlerFunc(fn)
}

func composeData(headers []byte, body []byte) []byte {
	d := make([]byte, 0)
	d = append(d, headers...)
	d = append(d, 0)
	d = append(d, body...)
	return d
}

func decomposeAndWriteData(key string, data []byte, w http.ResponseWriter) {
	log.Printf("[DEBUG] Cache found for key: %s", key)
	sp := findSplitter(data, 0)
	if sp > -1 {
		h := make(map[string][]string)
		if err := json.Unmarshal(data[:sp], &h); err != nil {
			log.Printf("[ERROR] %v", err)
		}
		for k, v := range h {
			w.Header()[k] = v
		}
		data = data[sp+1:]
	}
	w.Write(data)
}
