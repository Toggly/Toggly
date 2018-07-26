package cache

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

func TestCacheKey(t *testing.T) {
	assert := assert.New(t)
	url := "http://localhost/?a=1"
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(err)
	key := getKeyFromRequest(req)
	assert.Equal(key, url)
}

func mockBody() string {
	j, _ := json.Marshal(struct {
		A string `json:"a"`
	}{A: "test"})
	return string(j)
}

func mockFunction(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, struct {
		A string `json:"a"`
	}{A: "test"})
}

func TestCached(t *testing.T) {
	assert := assert.New(t)
	cache, _ := NewHashMapCache(true)
	cfn := Cached(mockFunction, cache)
	req, _ := http.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	b := mockBody()
	cfn(w, req)
	assert.Equal(b, strings.TrimSpace(w.Body.String()))

}
