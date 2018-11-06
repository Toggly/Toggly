package rest_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/server/rest"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

func TestCacheKey(t *testing.T) {
	assert := assert.New(t)
	url := "http://localhost/?a=1"
	req, err := http.NewRequest("GET", url, nil)

	c := context.WithValue(context.Background(), rest.CtxValueOwner, "ow1")
	req = req.WithContext(c)

	assert.Nil(err)
	key := rest.GetKeyFromRequest(req)
	assert.Equal("ow1::"+url, key)
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
	cache, _ := cache.NewHashMapCache()
	cfn := rest.Cached(mockFunction, cache)
	req, _ := http.NewRequest("GET", "http://localhost", nil)

	c := context.WithValue(context.Background(), rest.CtxValueOwner, "ow1")
	req = req.WithContext(c)

	w := httptest.NewRecorder()
	b := mockBody()
	cfn(w, req)
	assert.Equal(b, strings.TrimSpace(w.Body.String()))
	//TODO make more proper testing
}
