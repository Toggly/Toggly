package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/pkg/ctx"
)

//GetKeyFromRequest composes the key based on owner id and url
func GetKeyFromRequest(r *http.Request) string {
	return ctx.CtxOwner(r) + "::" + r.URL.String()
}

// Cached implements http.HandlerFunc caching
func Cached(next http.HandlerFunc, cache cache.DataCache) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if !cache.Enabled() {
			next.ServeHTTP(w, r)
			return
		}
		key := GetKeyFromRequest(r)

		data, err := cache.Get(key)
		if err != nil {
			log.Printf("[ERROR] Cache error: %v", err)
			data = nil
		}

		if data != nil {
			log.Printf("[DEBUG] 🔥 From cache: %s", key)
			decomposeAndWriteData(key, data, w)
			return
		}

		log.Printf("[DEBUG] 📀 From DB: %s", key)

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
