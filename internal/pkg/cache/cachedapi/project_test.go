package cachedapi_test

import (
	"encoding/json"
	"testing"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/api"
	asserts "github.com/stretchr/testify/assert"
)

func TestProjectCaching(t *testing.T) {
	assert := asserts.New(t)

	BeforeTest()

	engine, cache := getEngineAndCache()
	eng := engine.ForOwner("ow1").Projects()

	b, err := cache.Get("/own/ow1/project")
	assert.Nil(err)
	assert.Nil(b)

	eng.Create(&api.ProjectInfo{Code: "project1", Description: "Description 1"})
	b, err = cache.Get("/own/ow1/project")
	assert.Nil(err)
	assert.Nil(b)

	eng.List()

	var list []*domain.Project

	b, err = cache.Get("/own/ow1/project")
	assert.Nil(err)
	assert.NotNil(b)
	json.Unmarshal(b, &list)
	assert.Len(list, 1)

	b, err = cache.Get("/own/ow1/project/project1")
	assert.Nil(err)
	assert.Nil(b)

	_, err = eng.List()
	assert.Nil(err)

	proj := &domain.Project{}

	eng.Get("project1")
	b, err = cache.Get("/own/ow1/project/project1")
	assert.Nil(err)
	assert.NotNil(b)
	json.Unmarshal(b, proj)
	assert.Equal(domain.ProjectCode("project1"), proj.Code)
	assert.Equal("Description 1", proj.Description)

	eng.Update(&api.ProjectInfo{Code: "project1", Description: "Description 2"})
	b, err = cache.Get("/own/ow1/project")
	assert.Nil(err)
	assert.Nil(b)
	b, err = cache.Get("/own/ow1/project/project1")
	assert.Nil(err)
	assert.Nil(b)

	eng.Create(&api.ProjectInfo{Code: "project1", Description: "Description 1"})
	eng.List()
	eng.Get("project1")
	b, err = cache.Get("/own/ow1/project")
	assert.Nil(err)
	assert.NotNil(b)
	b, err = cache.Get("/own/ow1/project/project1")
	assert.Nil(err)
	assert.NotNil(b)

	eng.Delete("project1")
	b, err = cache.Get("/own/ow1/project")
	assert.Nil(err)
	assert.Nil(b)
	b, err = cache.Get("/own/ow1/project/project1")
	assert.Nil(err)
	assert.Nil(b)

	AfterTest()
}
